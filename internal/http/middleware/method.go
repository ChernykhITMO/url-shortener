package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
)

const methodNotAllowed = "method not allowed"

func RequireMethodMiddleware(method string, log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				log.Warn(methodNotAllowed,
					slog.String("method", r.Method),
					slog.String("expected_method", method),
					slog.String("path", r.URL.Path),
				)

				w.Header().Set("Allow", method)

				if err := respond.WriteJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed); err != nil {
					log.Error("write error response failed", slog.Any("err", err))
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
