package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Fabric struct {
	ID           string          `json:"id"`
	OwnerID      string          `json:"ownerId"`
	Name         string          `json:"name"`
	Material     string          `json:"material"`
	Pattern      string          `json:"pattern"`
	StyleTags    []string        `json:"styleTags"`
	Colors       json.RawMessage `json:"colors"`
	Memo         string          `json:"memo"`
	CoverAssetID string          `json:"coverAssetId,omitempty"`
	CoverURL     string          `json:"coverUrl,omitempty"`
	Pieces       []FabricPiece   `json:"pieces"`
	Version      int64           `json:"version"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type FabricPiece struct {
	ID              string     `json:"id"`
	FabricID        string     `json:"fabricId"`
	Label           string     `json:"label"`
	LengthMM        *int       `json:"lengthMm,omitempty"`
	WidthMM         *int       `json:"widthMm,omitempty"`
	Quantity        float64    `json:"quantity"`
	Unit            string     `json:"unit"`
	StorageLocation string     `json:"storageLocation"`
	PurchaseChannel string     `json:"purchaseChannel"`
	PurchaseDate    *time.Time `json:"purchaseDate,omitempty"`
	TotalPriceCent  *int64     `json:"totalPriceCent,omitempty"`
	Status          string     `json:"status"`
	Memo            string     `json:"memo"`
}

type Asset struct {
	ID         string `json:"id"`
	StorageKey string `json:"storageKey"`
	URL        string `json:"url"`
}

type CreateFabricInput struct {
	Name         string
	Material     string
	Pattern      string
	StyleTags    []string
	Colors       json.RawMessage
	Memo         string
	CoverAssetID string
	Piece        *CreateFabricPieceInput
}

type CreateFabricPieceInput struct {
	Label           string
	LengthMM        *int
	WidthMM         *int
	StorageLocation string
	PurchaseChannel string
	PurchaseDate    *time.Time
	TotalPriceCent  *int64
	Memo            string
}

type Project struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Status    string    `json:"status"`
	Memo      string    `json:"memo"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store struct {
	pool              *pgxpool.Pool
	ownerID           string
	storagePublicBase string
}

func NewStore(ctx context.Context, pool *pgxpool.Pool, storagePublicBase string) (*Store, error) {
	var ownerID string
	err := pool.QueryRow(ctx, `
		INSERT INTO users (external_provider, external_id, nickname)
		VALUES ('system', 'default', '布迹用户')
		ON CONFLICT (external_provider, external_id)
		DO UPDATE SET updated_at = now()
		RETURNING id::text
	`).Scan(&ownerID)
	if err != nil {
		return nil, fmt.Errorf("ensure default user: %w", err)
	}
	return &Store{pool: pool, ownerID: ownerID, storagePublicBase: storagePublicBase}, nil
}

func (s *Store) ListFabrics(ctx context.Context) ([]Fabric, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT f.id::text, f.owner_id::text, f.name, f.material, f.pattern, f.style_tags,
		       f.colors, f.memo, COALESCE(f.cover_asset_id::text, ''), COALESCE(a.storage_key, ''),
		       f.version, f.created_at, f.updated_at
		FROM fabrics f
		LEFT JOIN assets a ON a.id = f.cover_asset_id
		WHERE f.owner_id = $1 AND f.deleted_at IS NULL
		ORDER BY f.updated_at DESC
	`, s.ownerID)
	if err != nil {
		return nil, fmt.Errorf("list fabrics: %w", err)
	}
	defer rows.Close()

	items := make([]Fabric, 0)
	for rows.Next() {
		var item Fabric
		var storageKey string
		if err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Material, &item.Pattern, &item.StyleTags, &item.Colors, &item.Memo, &item.CoverAssetID, &storageKey, &item.Version, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fabric: %w", err)
		}
		if storageKey != "" && s.storagePublicBase != "" {
			item.CoverURL = s.storagePublicBase + "/" + storageKey
		}
		item.Pieces, err = s.listPieces(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateFabric(ctx context.Context, input CreateFabricInput) (Fabric, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Fabric{}, fmt.Errorf("begin create fabric: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var item Fabric
	colors := input.Colors
	if len(colors) == 0 {
		colors = json.RawMessage(`[]`)
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO fabrics (owner_id, name, material, pattern, style_tags, colors, memo, cover_asset_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '')::uuid)
		RETURNING id::text, owner_id::text, name, material, pattern, style_tags,
		          colors, memo, COALESCE(cover_asset_id::text, ''), version, created_at, updated_at
	`, s.ownerID, input.Name, input.Material, input.Pattern, input.StyleTags, colors, input.Memo, input.CoverAssetID).Scan(
		&item.ID, &item.OwnerID, &item.Name, &item.Material, &item.Pattern,
		&item.StyleTags, &item.Colors, &item.Memo, &item.CoverAssetID, &item.Version, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return Fabric{}, fmt.Errorf("create fabric: %w", err)
	}
	if input.Piece != nil {
		p := input.Piece
		if p.Label == "" {
			p.Label = "新入布块"
		}
		var piece FabricPiece
		err = tx.QueryRow(ctx, `INSERT INTO fabric_pieces
			(owner_id, fabric_id, label, length_mm, width_mm, storage_location, purchase_channel, purchase_date, total_price_cent, memo)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			RETURNING id::text, fabric_id::text, label, length_mm, width_mm, quantity::float8, unit,
			storage_location, purchase_channel, purchase_date, total_price_cent, status::text, memo`,
			s.ownerID, item.ID, p.Label, p.LengthMM, p.WidthMM, p.StorageLocation, p.PurchaseChannel, p.PurchaseDate, p.TotalPriceCent, p.Memo,
		).Scan(&piece.ID, &piece.FabricID, &piece.Label, &piece.LengthMM, &piece.WidthMM, &piece.Quantity, &piece.Unit, &piece.StorageLocation, &piece.PurchaseChannel, &piece.PurchaseDate, &piece.TotalPriceCent, &piece.Status, &piece.Memo)
		if err != nil {
			return Fabric{}, fmt.Errorf("create fabric piece: %w", err)
		}
		item.Pieces = []FabricPiece{piece}
	}
	if err := tx.Commit(ctx); err != nil {
		return Fabric{}, fmt.Errorf("commit create fabric: %w", err)
	}
	return item, nil
}

func (s *Store) SaveAsset(ctx context.Context, storageKey, mimeType string, width, height int, size int64, sha256 string) (Asset, error) {
	var asset Asset
	err := s.pool.QueryRow(ctx, `INSERT INTO assets
		(owner_id, storage_key, mime_type, width, height, size_bytes, sha256)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (owner_id, sha256) DO UPDATE SET updated_at = now()
		RETURNING id::text, storage_key`, s.ownerID, storageKey, mimeType, width, height, size, sha256).Scan(&asset.ID, &asset.StorageKey)
	if err != nil {
		return Asset{}, fmt.Errorf("save asset: %w", err)
	}
	if s.storagePublicBase != "" {
		asset.URL = s.storagePublicBase + "/" + asset.StorageKey
	}
	return asset, nil
}

func (s *Store) OwnerID() string { return s.ownerID }

func (s *Store) listPieces(ctx context.Context, fabricID string) ([]FabricPiece, error) {
	rows, err := s.pool.Query(ctx, `SELECT id::text, fabric_id::text, label, length_mm, width_mm,
		quantity::float8, unit, storage_location, purchase_channel, purchase_date,
		total_price_cent, status::text, memo FROM fabric_pieces
		WHERE fabric_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC`, fabricID)
	if err != nil {
		return nil, fmt.Errorf("list fabric pieces: %w", err)
	}
	defer rows.Close()
	items := make([]FabricPiece, 0)
	for rows.Next() {
		var p FabricPiece
		if err := rows.Scan(&p.ID, &p.FabricID, &p.Label, &p.LengthMM, &p.WidthMM, &p.Quantity, &p.Unit, &p.StorageLocation, &p.PurchaseChannel, &p.PurchaseDate, &p.TotalPriceCent, &p.Status, &p.Memo); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (s *Store) ListProjects(ctx context.Context) ([]Project, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id::text, owner_id::text, name, kind, status::text, memo,
		       version, created_at, updated_at
		FROM projects
		WHERE owner_id = $1 AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`, s.ownerID)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	items := make([]Project, 0)
	for rows.Next() {
		var item Project
		if err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Kind, &item.Status, &item.Memo, &item.Version, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateProject(ctx context.Context, name, kind, memo string) (Project, error) {
	var item Project
	err := s.pool.QueryRow(ctx, `
		INSERT INTO projects (owner_id, name, kind, memo)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, owner_id::text, name, kind, status::text, memo,
		          version, created_at, updated_at
	`, s.ownerID, name, kind, memo).Scan(
		&item.ID, &item.OwnerID, &item.Name, &item.Kind, &item.Status,
		&item.Memo, &item.Version, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return Project{}, fmt.Errorf("create project: %w", err)
	}
	return item, nil
}
