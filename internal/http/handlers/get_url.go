package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/dto"
	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")

	url, err := h.service.GetURL(r.Context(), alias)
	if err != nil {
		h.log.Error("get url failed", slog.Any("err", err), slog.String("alias", alias))
		h.writeServiceError(w, err)
		return
	}

	resp := dto.GetURLResponse{URL: url}
	if err := respond.WriteJSON(w, http.StatusOK, resp); err != nil {
		h.log.Error("write success response failed", slog.Any("err", err))
	}
}
