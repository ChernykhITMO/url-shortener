package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ChernykhITMO/url-shortener/internal/http/dto"
	"github.com/ChernykhITMO/url-shortener/internal/http/respond"
)

const maxCreateAliasBodyBytes = 1024

func (h *Handler) CreateAlias(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	limited := http.MaxBytesReader(w, r.Body, maxCreateAliasBodyBytes)
	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields()

	var req dto.CreateAliasRequest
	if err := dec.Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			h.writeJSONError(w, http.StatusRequestEntityTooLarge, msgPayloadTooLarge)
			return
		}
		h.log.Error("decode failed", slog.Any("err", err))
		h.writeJSONError(w, http.StatusBadRequest, msgInvalidJSON)
		return
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		h.log.Error("unexpected data after JSON object", slog.Any("err", err))
		h.writeJSONError(w, http.StatusBadRequest, msgInvalidJSON)
		return
	}

	alias, err := h.service.CreateAlias(r.Context(), req.URL)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	resp := dto.CreateAliasResponse{Alias: alias}
	if err := respond.WriteJSON(w, http.StatusCreated, resp); err != nil {
		h.log.Error("write success response failed", slog.Any("err", err))
	}
}
