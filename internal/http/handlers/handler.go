package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
	"github.com/ChernykhITMO/url-shortener/internal/services"
)

const (
	msgInvalidJSON     = "invalid json"
	msgInvalidURL      = "invalid url"
	msgNotFound        = "not found"
	msgTryAgainLater   = "try again later"
	msgInternalError   = "internal error"
	msgPayloadTooLarge = "payload too large"
	msgInvalidAlias    = "invalid alias"
)

type Service interface {
	CreateAlias(ctx context.Context, originalURL string) (string, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

type Handler struct {
	log     *slog.Logger
	Service Service
}

func New(log *slog.Logger, service Service) *Handler {
	return &Handler{
		log:     log,
		Service: service,
	}
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidURL):
		if writeErr := respond.WriteJSONError(w, http.StatusBadRequest, msgInvalidURL); writeErr != nil {
			h.log.Error("write error response failed", slog.Any("err", writeErr))
		}
	case errors.Is(err, services.ErrNotFound):
		if writeErr := respond.WriteJSONError(w, http.StatusNotFound, msgNotFound); writeErr != nil {
			h.log.Error("write error response failed", slog.Any("err", writeErr))
		}
	case errors.Is(err, services.ErrAttemptsExceeded):
		if writeErr := respond.WriteJSONError(w, http.StatusServiceUnavailable, msgTryAgainLater); writeErr != nil {
			h.log.Error("write error response failed", slog.Any("err", writeErr))
		}
	default:
		if writeErr := respond.WriteJSONError(w, http.StatusInternalServerError, msgInternalError); writeErr != nil {
			h.log.Error("write error response failed", slog.Any("err", writeErr))
		}
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	if err := respond.WriteJSONError(w, status, msg); err != nil {
		h.log.Error("write error response failed", slog.Any("err", err))
	}
}
