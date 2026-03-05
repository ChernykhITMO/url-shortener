package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
)

const msgInternalError = "internal error"

type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.wroteHeader = true
	return rw.ResponseWriter.Write(b)
}

func RecoverMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respWr := &responseWriter{ResponseWriter: w}
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)

					if !respWr.wroteHeader {
						if err := respond.WriteJSONError(respWr, http.StatusInternalServerError, msgInternalError); err != nil {
							log.Error("write error response failed", slog.Any("err", err))
						}
					}
				}
			}()
			next.ServeHTTP(respWr, r)
		})
	}
}
