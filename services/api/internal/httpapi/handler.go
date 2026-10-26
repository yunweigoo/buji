package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/yunweigoo/buji/services/api/internal/config"
)

type response struct {
	Data      any    `json:"data"`
	RequestID string `json:"requestId"`
}

func NewHandler(logger *slog.Logger, cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{
			Data: map[string]any{
				"status":      "ok",
				"service":     "buji-api",
				"environment": cfg.Environment,
				"time":        time.Now().UTC().Format(time.RFC3339),
			},
			RequestID: requestID(r),
		})
	})
	mux.HandleFunc("GET /v1/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{
			Data: map[string]any{
				"fabrics":  []any{},
				"projects": []any{},
				"features": map[string]bool{
					"cloudSync":        cfg.DatabaseURL != "",
					"feishuBackoffice": cfg.FeishuAppID != "" && cfg.FeishuBaseToken != "",
				},
			},
			RequestID: requestID(r),
		})
	})

	return recoverMiddleware(logger, accessLogMiddleware(logger, corsMiddleware(mux)))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
