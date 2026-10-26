package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Fabric struct {
	ID        string          `json:"id"`
	OwnerID   string          `json:"ownerId"`
	Name      string          `json:"name"`
	Material  string          `json:"material"`
	Pattern   string          `json:"pattern"`
	StyleTags []string        `json:"styleTags"`
	Colors    json.RawMessage `json:"colors"`
	Memo      string          `json:"memo"`
	Version   int64           `json:"version"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
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
	pool    *pgxpool.Pool
	ownerID string
}

func NewStore(ctx context.Context, pool *pgxpool.Pool) (*Store, error) {
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
	return &Store{pool: pool, ownerID: ownerID}, nil
}

func (s *Store) ListFabrics(ctx context.Context) ([]Fabric, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id::text, owner_id::text, name, material, pattern, style_tags,
		       colors, memo, version, created_at, updated_at
		FROM fabrics
		WHERE owner_id = $1 AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`, s.ownerID)
	if err != nil {
		return nil, fmt.Errorf("list fabrics: %w", err)
	}
	defer rows.Close()

	items := make([]Fabric, 0)
	for rows.Next() {
		var item Fabric
		if err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Material, &item.Pattern, &item.StyleTags, &item.Colors, &item.Memo, &item.Version, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fabric: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateFabric(ctx context.Context, name, material, pattern, memo string) (Fabric, error) {
	var item Fabric
	err := s.pool.QueryRow(ctx, `
		INSERT INTO fabrics (owner_id, name, material, pattern, memo)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, owner_id::text, name, material, pattern, style_tags,
		          colors, memo, version, created_at, updated_at
	`, s.ownerID, name, material, pattern, memo).Scan(
		&item.ID, &item.OwnerID, &item.Name, &item.Material, &item.Pattern,
		&item.StyleTags, &item.Colors, &item.Memo, &item.Version, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return Fabric{}, fmt.Errorf("create fabric: %w", err)
	}
	return item, nil
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
