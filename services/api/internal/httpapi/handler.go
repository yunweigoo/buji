package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
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
	CreateFabric(context.Context, string, string, string, string) (database.Fabric, error)
	ListProjects(context.Context) ([]database.Project, error)
	CreateProject(context.Context, string, string, string) (database.Project, error)
}

func NewHandler(logger *slog.Logger, cfg config.Config, pool *pgxpool.Pool, store dataStore) http.Handler {
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
			Name     string `json:"name"`
			Material string `json:"material"`
			Pattern  string `json:"pattern"`
			Memo     string `json:"memo"`
		}
		if err := decodeJSON(r, &input); err != nil || input.Name == "" {
			writeAPIError(logger, w, r, http.StatusBadRequest, "invalid_request", err)
			return
		}
		item, err := store.CreateFabric(r.Context(), input.Name, input.Material, input.Pattern, input.Memo)
		if err != nil {
			writeAPIError(logger, w, r, http.StatusInternalServerError, "database_error", err)
			return
		}
		writeJSON(w, http.StatusCreated, response{Data: item, RequestID: requestID(r)})
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
