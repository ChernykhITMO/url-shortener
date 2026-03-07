package services

import (
	"context"
	"errors"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "services.GetURL"

	if err := validateAlias(alias); err != nil {
		return "", wrapError(op, err)
	}

	url, err := s.storage.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", wrapError(op, ErrNotFound)
		}
		return "", wrapError(op, err)
	}

	return url, nil
}
