package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yunweigoo/buji/services/api/internal/config"
	"github.com/yunweigoo/buji/services/api/internal/database"
	"github.com/yunweigoo/buji/services/api/internal/httpapi"
	"github.com/yunweigoo/buji/services/api/internal/storage"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if cfg.DatabaseURL == "" {
		logger.Error("database URL is required", "accepted_variables", "DATABASE_URL, POSTGRES_URL, POSTGRES_URL_NON_POOLING")
		os.Exit(1)
	}

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer startupCancel()
	pool, err := database.Open(startupCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := database.Migrate(startupCtx, pool); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}
	objects := storage.NewSupabase(cfg.SupabaseURL, cfg.SupabaseKey)
	if !objects.Ready() {
		logger.Error("Supabase Storage configuration is required", "accepted_variables", "SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY or SUPABASE_SECRET_KEY")
		os.Exit(1)
	}
	if err := objects.EnsureBucket(startupCtx); err != nil {
		logger.Error("Supabase Storage initialization failed", "error", err)
		os.Exit(1)
	}
	storagePublicBase := strings.TrimRight(cfg.SupabaseURL, "/") + "/storage/v1/object/public/" + storage.Bucket
	store, err := database.NewStore(startupCtx, pool, storagePublicBase)
	if err != nil {
		logger.Error("database store initialization failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database ready", "migration_version", 1)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewHandler(logger, cfg, pool, store, objects),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("api listening", "port", cfg.Port, "environment", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
