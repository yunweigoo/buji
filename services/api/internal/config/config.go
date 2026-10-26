package config

import "os"

type Config struct {
	Port            string
	Environment     string
	DatabaseURL     string
	CloudBaseEnvID  string
	FeishuAppID     string
	FeishuAppSecret string
	FeishuBaseToken string
}

func Load() Config {
	return Config{
		Port:            envOr("PORT", "8080"),
		Environment:     envOr("APP_ENV", "development"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		CloudBaseEnvID:  os.Getenv("CLOUDBASE_ENV_ID"),
		FeishuAppID:     os.Getenv("FEISHU_APP_ID"),
		FeishuAppSecret: os.Getenv("FEISHU_APP_SECRET"),
		FeishuBaseToken: os.Getenv("FEISHU_BASE_APP_TOKEN"),
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
