package router

import (
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/handlers"
	"github.com/ChernykhITMO/url-shortener/internal/http/middleware"
)

func New(handler *handlers.Handler, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	createAliasHandler := middleware.RecoverMiddleware(log)(
		middleware.LoggingMiddleware(log)(
			middleware.RequireMethodMiddleware(http.MethodPost, log)(http.HandlerFunc(handler.CreateAlias)),
		),
	)
	getURLHandler := middleware.RecoverMiddleware(log)(
		middleware.LoggingMiddleware(log)(
			middleware.RequireMethodMiddleware(http.MethodGet, log)(http.HandlerFunc(handler.GetURL)),
		),
	)

	mux.Handle("/url", createAliasHandler)
	mux.Handle("/url/{alias}", getURLHandler)

	return mux
}
