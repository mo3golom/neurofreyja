package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken       string
	BotUsername         string
	OpenRouterAPIKey    string
	OpenRouterBaseURL   string
	OpenRouterModelTit  string
	OpenRouterModelDesc string
	PGDSN               string
	S3Endpoint          string
	S3AccessKey         string
	S3SecretKey         string
	S3Bucket            string
	S3Region            string
	S3PathStyle         bool
	DeleteAfterMinutes  int
	DeleteSweepInterval time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{}
	cfg.TelegramToken = strings.TrimSpace(os.Getenv("TELEGRAM_TOKEN"))
	cfg.BotUsername = strings.TrimSpace(os.Getenv("BOT_USERNAME"))
	cfg.OpenRouterAPIKey = strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))
	cfg.OpenRouterBaseURL = strings.TrimSpace(os.Getenv("OPENROUTER_BASE_URL"))
	if cfg.OpenRouterBaseURL == "" {
		cfg.OpenRouterBaseURL = "https://openrouter.ai/api/v1"
	}
	cfg.OpenRouterModelTit = strings.TrimSpace(os.Getenv("OPENROUTER_MODEL_TITLES"))
	if cfg.OpenRouterModelTit == "" {
		cfg.OpenRouterModelTit = "openai/gpt-4.1"
	}
	cfg.OpenRouterModelDesc = strings.TrimSpace(os.Getenv("OPENROUTER_MODEL_DESC"))
	if cfg.OpenRouterModelDesc == "" {
		cfg.OpenRouterModelDesc = "openai/gpt-4.1"
	}

	cfg.PGDSN = strings.TrimSpace(os.Getenv("PG_DSN"))

	cfg.S3Endpoint = strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	cfg.S3AccessKey = strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	cfg.S3SecretKey = strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))
	cfg.S3Bucket = strings.TrimSpace(os.Getenv("S3_BUCKET"))
	if cfg.S3Bucket == "" {
		cfg.S3Bucket = "neurofreyja"
	}
	cfg.S3Region = strings.TrimSpace(os.Getenv("S3_REGION"))

	cfg.S3PathStyle = true
	if raw := strings.TrimSpace(os.Getenv("S3_PATH_STYLE")); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			cfg.S3PathStyle = v
		}
	}

	cfg.DeleteAfterMinutes = 10
	if raw := strings.TrimSpace(os.Getenv("DELETE_AFTER_MINUTES")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			cfg.DeleteAfterMinutes = v
		}
	}

	cfg.DeleteSweepInterval = time.Minute
	if raw := strings.TrimSpace(os.Getenv("DELETE_SWEEP_INTERVAL")); raw != "" {
		if v, err := time.ParseDuration(raw); err == nil {
			cfg.DeleteSweepInterval = v
		}
	}

	return cfg
}
