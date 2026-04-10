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
	msgAliasTaken      = "alias already taken"
	msgUnsupportedType = "content type must be application/json"
)

type Service interface {
	CreateAlias(ctx context.Context, originalURL, requestedAlias string) (string, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

type Handler struct {
	log     *slog.Logger
	service Service
}

func New(log *slog.Logger, service Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidURL):
		h.writeJSONError(w, http.StatusBadRequest, msgInvalidURL)
	case errors.Is(err, services.ErrNotFound):
		h.writeJSONError(w, http.StatusNotFound, msgNotFound)
	case errors.Is(err, services.ErrAttemptsExceeded):
		h.writeJSONError(w, http.StatusServiceUnavailable, msgTryAgainLater)
	case errors.Is(err, services.ErrInvalidAlias):
		h.writeJSONError(w, http.StatusBadRequest, msgInvalidAlias)
	case errors.Is(err, services.ErrAliasTaken):
		h.writeJSONError(w, http.StatusConflict, msgAliasTaken)
	default:
		h.writeJSONError(w, http.StatusInternalServerError, msgInternalError)
	}
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, msg string) {
	if err := respond.WriteJSONError(w, status, msg); err != nil {
		h.log.Error("write error response failed", slog.Any("err", err))
	}
}
