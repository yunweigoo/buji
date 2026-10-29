package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Environment     string
	DatabaseURL     string
	SupabaseURL     string
	SupabaseKey     string
	FeishuAppID     string
	FeishuAppSecret string
	FeishuBaseToken string
}

func Load() Config {
	// Environment variables supplied by the process always win. These files are
	// only a convenience for local development from either the repository root
	// or services/api.
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../../.env")

	return Config{
		Port:            envOr("PORT", "8080"),
		Environment:     envOr("APP_ENV", "development"),
		DatabaseURL:     firstEnv("DATABASE_URL", "POSTGRES_URL", "POSTGRES_URL_NON_POOLING"),
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		SupabaseKey:     firstEnv("SUPABASE_SERVICE_ROLE_KEY", "SUPABASE_SECRET_KEY"),
		FeishuAppID:     os.Getenv("FEISHU_APP_ID"),
		FeishuAppSecret: os.Getenv("FEISHU_APP_SECRET"),
		FeishuBaseToken: os.Getenv("FEISHU_BASE_APP_TOKEN"),
	}
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
