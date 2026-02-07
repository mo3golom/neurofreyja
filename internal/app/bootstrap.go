package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"neurofreyja/internal/entities/card"
	"neurofreyja/internal/entities/drawn"
	"neurofreyja/internal/entities/history"
	"neurofreyja/internal/shared/config"
	"neurofreyja/internal/shared/db"
	"neurofreyja/internal/shared/llm"
	"neurofreyja/internal/shared/logger"
	"neurofreyja/internal/shared/storage"
	"neurofreyja/internal/shared/telegram"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type App struct {
	Config    *config.Config
	Logger    *logrus.Logger
	DB        *sqlx.DB
	Bot       *telebot.Bot
	Messenger *telegram.Messenger
	Cards     card.Repository
	Drawn     drawn.Repository
	History   history.Repository
	Storage   *storage.Client
	LLM       *llm.Client
}

func Bootstrap() (*App, error) {
	cfg := config.Load()
	log := logger.New()

	if cfg.TelegramToken == "" {
		return nil, errors.New("TELEGRAM_TOKEN is required")
	}
	if cfg.PGDSN == "" {
		return nil, errors.New("PG_DSN is required")
	}

	dbConn, err := db.Connect(cfg.PGDSN)
	if err != nil {
		return nil, err
	}

	historyRepo := history.NewRepoSQLX(dbConn)
	cardRepo := card.NewRepoSQLX(dbConn)
	drawnRepo := drawn.NewRepoSQLX(dbConn)

	storageClient, err := storage.New(context.Background(), storage.Config{
		Endpoint:  cfg.S3Endpoint,
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Bucket:    cfg.S3Bucket,
		Region:    cfg.S3Region,
		PathStyle: cfg.S3PathStyle,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramToken,
		Client: httpClient,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	if cfg.BotUsername == "" {
		cfg.BotUsername = bot.Me.Username
	}

	messenger := telegram.NewMessenger(bot, log, historyRepo)
	llmClient := llm.New(cfg.OpenRouterBaseURL, cfg.OpenRouterAPIKey)

	return &App{
		Config:    cfg,
		Logger:    log,
		DB:        dbConn,
		Bot:       bot,
		Messenger: messenger,
		Cards:     cardRepo,
		Drawn:     drawnRepo,
		History:   historyRepo,
		Storage:   storageClient,
		LLM:       llmClient,
	}, nil
}
