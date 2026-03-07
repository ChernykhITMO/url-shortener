package main

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	app_http "github.com/ChernykhITMO/url-shortener/cmd/app/http"
	"github.com/ChernykhITMO/url-shortener/internal/config"
	"github.com/ChernykhITMO/url-shortener/internal/services"
	"github.com/ChernykhITMO/url-shortener/internal/storage/inmemory"
	"github.com/ChernykhITMO/url-shortener/internal/storage/postgres"
)

const (
	inMemoryStorage = "inmemory"
	postgresStorage = "postgres"

	envLocal = "local"
)

func main() {
	run()
}

func run() {
	storageStr := flag.String("storage", inMemoryStorage, "storage type")
	configPath := flag.String("config", "./config/local.yaml", "path to config yaml")
	flag.Parse()

	bootstrapLog := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		bootstrapLog.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := setupLogger(cfg.Env)

	var storage services.Storage

	switch *storageStr {
	case inMemoryStorage:
		storage = inmemory.New()
	case postgresStorage:
		if cfg.Postgres.DSN == "" {
			log.Error("postgres.dsn is required for postgres storage")
			os.Exit(1)
		}
		storage, err = postgres.New(ctx, cfg.Postgres)
		if err != nil {
			log.Error("failed create storage", slog.Any("err", err))
			os.Exit(1)
		}
	default:
		log.Error("unknown storage type", slog.String("storage", *storageStr))
		os.Exit(1)
	}

	if c, ok := storage.(io.Closer); ok {
		defer func() {
			if err := c.Close(); err != nil {
				log.Error("failed to close storage", slog.Any("err", err))
			}
		}()
	}

	if err := app_http.Run(ctx, storage, cfg, log); err != nil {
		log.Error("application stopped", slog.Any("err", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}
	return log
}
