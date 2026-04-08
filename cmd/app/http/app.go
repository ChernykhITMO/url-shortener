package app_http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/config"
	"github.com/ChernykhITMO/url-shortener/internal/http/handlers"
	"github.com/ChernykhITMO/url-shortener/internal/http/router"
	"github.com/ChernykhITMO/url-shortener/internal/services"
)

func Run(ctx context.Context, storage services.Storage, cfg *config.Config, log *slog.Logger) error {
	const op = "app_http.Run"

	log = log.With(slog.String("op", op))
	log.Info("http server config",
		slog.String("addr", cfg.HTTPServer.Address),
		slog.Duration("read_header_timeout", cfg.HTTPServer.ReadHeaderTimeout),
		slog.Duration("read_timeout", cfg.HTTPServer.ReadTimeout),
		slog.Duration("write_timeout", cfg.HTTPServer.WriteTimeout),
		slog.Duration("idle_timeout", cfg.HTTPServer.IdleTimeout),
		slog.Duration("shutdown_timeout", cfg.HTTPServer.ShutdownTimeout),
	)

	log.Info("starting application...")

	service := services.New(storage, cfg.Service.MaxAttempts, cfg.Service.AliasLength)
	handler := handlers.New(log, service, cfg.HTTPServer.MaxBodyBytes)
	mux := router.New(handler, log)

	srv := &http.Server{
		Addr:              cfg.HTTPServer.Address,
		Handler:           mux,
		ReadHeaderTimeout: cfg.HTTPServer.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      cfg.HTTPServer.WriteTimeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("received shutdown signal")
	case err := <-errCh:
		log.Error("server stopped", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	shutdownCtx, stop := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer stop()

	if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
