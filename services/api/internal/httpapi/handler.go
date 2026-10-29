package httpapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunweigoo/buji/services/api/internal/config"
	"github.com/yunweigoo/buji/services/api/internal/database"
)

type response struct {
	Data      any    `json:"data"`
	RequestID string `json:"requestId"`
}

type dataStore interface {
	ListFabrics(context.Context) ([]database.Fabric, error)
	CreateFabric(context.Context, database.CreateFabricInput) (database.Fabric, error)
	SaveAsset(context.Context, string, string, int, int, int64, string) (database.Asset, error)
	OwnerID() string
	ListProjects(context.Context) ([]database.Project, error)
	CreateProject(context.Context, string, string, string) (database.Project, error)
}

type objectStorage interface {
	Ready() bool
	Upload(context.Context, string, string, string, []byte) (string, error)
}

func NewHandler(logger *slog.Logger, cfg config.Config, pool *pgxpool.Pool, store dataStore, objects objectStorage) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		databaseStatus := "not_configured"
		serviceStatus := "degraded"
		status := http.StatusServiceUnavailable
		if pool != nil {
			if err := pool.Ping(r.Context()); err == nil {
				databaseStatus = "ok"
				serviceStatus = "ok"
				status = http.StatusOK
			} else {
				databaseStatus = "unavailable"
			}
		}
		writeJSON(w, status, response{
			Data: map[string]any{
				"status":      serviceStatus,
				"database":    databaseStatus,
				"service":     "buji-api",
				"environment": cfg.Environment,
				"time":        time.Now().UTC().Format(time.RFC3339),
			},
			RequestID: requestID(r),
		})
	})
	mux.HandleFunc("GET /v1/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		fabrics, err := store.ListFabrics(r.Context())
		if err != nil {
			writeAPIError(logger, w, r, http.StatusInternalServerError, "database_error", err)
			return
		}
		projects, err := store.ListProjects(r.Context())
		if err != nil {
			writeAPIError(logger, w, r, http.StatusInternalServerError, "database_error", err)
			return
		}
		writeJSON(w, http.StatusOK, response{
			Data: map[string]any{
				"fabrics":  fabrics,
				"projects": projects,
				"features": map[string]bool{
					"cloudSync":        pool != nil,
					"feishuBackoffice": cfg.FeishuAppID != "" && cfg.FeishuBaseToken != "",
				},
			},
			RequestID: requestID(r),
		})
	})
	mux.HandleFunc("POST /v1/fabrics", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name         string          `json:"name"`
			Material     string          `json:"material"`
			Pattern      string          `json:"pattern"`
			StyleTags    []string        `json:"styleTags"`
			Colors       json.RawMessage `json:"colors"`
			Memo         string          `json:"memo"`
			CoverAssetID string          `json:"coverAssetId"`
			Piece        *struct {
				Label           string `json:"label"`
				LengthMM        *int   `json:"lengthMm"`
				WidthMM         *int   `json:"widthMm"`
				StorageLocation string `json:"storageLocation"`
				PurchaseChannel string `json:"purchaseChannel"`
				PurchaseDate    string `json:"purchaseDate"`
				TotalPriceCent  *int64 `json:"totalPriceCent"`
				Memo            string `json:"memo"`
			} `json:"piece"`
		}
		if err := decodeJSON(r, &input); err != nil || strings.TrimSpace(input.Name) == "" {
			writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_request", err)
			return
		}
		createInput := database.CreateFabricInput{
			Name: strings.TrimSpace(input.Name), Material: strings.TrimSpace(input.Material),
			Pattern: strings.TrimSpace(input.Pattern), StyleTags: input.StyleTags,
			Colors: input.Colors, Memo: strings.TrimSpace(input.Memo), CoverAssetID: input.CoverAssetID,
		}
		if input.Piece != nil {
			piece := &database.CreateFabricPieceInput{
				Label: input.Piece.Label, LengthMM: input.Piece.LengthMM, WidthMM: input.Piece.WidthMM,
				StorageLocation: input.Piece.StorageLocation, PurchaseChannel: input.Piece.PurchaseChannel,
				TotalPriceCent: input.Piece.TotalPriceCent, Memo: input.Piece.Memo,
			}
			if input.Piece.PurchaseDate != "" {
				parsed, err := time.Parse("2006-01-02", input.Piece.PurchaseDate)
				if err != nil {
					writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_purchase_date", err)
					return
				}
				piece.PurchaseDate = &parsed
			}
			createInput.Piece = piece
		}
		item, err := store.CreateFabric(r.Context(), createInput)
		if err != nil {
			writeAPIError(logger, w, r, http.StatusInternalServerError, "database_error", err)
			return
		}
		writeJSON(w, http.StatusCreated, response{Data: item, RequestID: requestID(r)})
	})
	mux.HandleFunc("POST /v1/assets", func(w http.ResponseWriter, r *http.Request) {
		if objects == nil || !objects.Ready() {
			writeAPIError(logger, w, r, http.StatusServiceUnavailable, "storage_not_configured", nil)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 21<<20)
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_upload", err)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			writeAPIError(logger, w, r, http.StatusBadRequest, "file_required", err)
			return
		}
		defer file.Close()
		asset, err := uploadImage(r.Context(), objects, store, file, header)
		if err != nil {
			writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_image", err)
			return
		}
		writeJSON(w, http.StatusCreated, response{Data: asset, RequestID: requestID(r)})
	})
	mux.HandleFunc("POST /v1/projects", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
			Kind string `json:"kind"`
			Memo string `json:"memo"`
		}
		if err := decodeJSON(r, &input); err != nil || input.Name == "" {
			writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_request", err)
			return
		}
		item, err := store.CreateProject(r.Context(), input.Name, input.Kind, input.Memo)
		if err != nil {
			writeAPIError(logger, w, r, http.StatusInternalServerError, "database_error", err)
			return
		}
		writeJSON(w, http.StatusCreated, response{Data: item, RequestID: requestID(r)})
	})

	return recoverMiddleware(logger, accessLogMiddleware(logger, corsMiddleware(mux)))
}

func uploadImage(ctx context.Context, objects objectStorage, store dataStore, file multipart.File, header *multipart.FileHeader) (database.Asset, error) {
	data, err := io.ReadAll(io.LimitReader(file, (20<<20)+1))
	if err != nil {
		return database.Asset{}, err
	}
	if len(data) == 0 || len(data) > 20<<20 {
		return database.Asset{}, errors.New("image must be between 1 byte and 20 MB")
	}
	contentType := http.DetectContentType(data)
	extension := map[string]string{"image/jpeg": ".jpg", "image/png": ".png"}[contentType]
	if extension == "" {
		return database.Asset{}, fmt.Errorf("unsupported image type %q", contentType)
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return database.Asset{}, fmt.Errorf("decode image %s: %w", filepath.Base(header.Filename), err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	key, err := objects.Upload(ctx, store.OwnerID(), extension, contentType, data)
	if err != nil {
		return database.Asset{}, err
	}
	return store.SaveAsset(ctx, key, contentType, config.Width, config.Height, int64(len(data)), digest)
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeAPIError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, code string, err error) {
	if err != nil {
		logger.Error("request failed", "code", code, "method", r.Method, "path", r.URL.Path, "error", err)
	}
	writeJSON(w, status, map[string]any{"error": code, "requestId": requestID(r)})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
