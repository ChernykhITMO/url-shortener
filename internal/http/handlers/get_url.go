package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
	"github.com/ChernykhITMO/url-shortener/internal/transport/http/dto"
)

const (
	lengthAlias = 10
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	if !validationAlias(alias) {
		h.writeError(w, http.StatusBadRequest, msgInvalidAlias)
		return
	}

	url, err := h.Service.GetURL(r.Context(), alias)
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

func validationAlias(s string) bool {
	if len(s) != lengthAlias {
		return false
	}

	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' ||
			s[i] >= 'A' && s[i] <= 'Z' ||
			s[i] >= '0' && s[i] <= '9' ||
			s[i] == '_' {
			continue
		}
		return false
	}
	return true
}
