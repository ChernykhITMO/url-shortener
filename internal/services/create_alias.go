package services

import (
	"context"
	"errors"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

var ErrAttemptsExceeded = errors.New("attempts exceeded")

func (s *Service) CreateAlias(ctx context.Context, originalURL string) (string, error) {
	const op = "services.CreateAlias"

	if err := validateURL(originalURL); err != nil {
		return "", wrapError(op, err)
	}

	alias, err := s.storage.GetAliasByURL(ctx, originalURL)
	if err == nil {
		return alias, nil
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return "", wrapError(op, err)
	}

	for i := 0; i < s.maxAttempts; i++ {
		alias, err = s.generateAlias()
		if err != nil {
			return "", wrapError(op, err)
		}

		err = s.storage.Create(ctx, alias, originalURL)
		if err == nil {
			return alias, nil
		}

		if errors.Is(err, storage.ErrAliasConflict) {
			continue
		}

		if errors.Is(err, storage.ErrURLConflict) {
			existing, err := s.storage.GetAliasByURL(ctx, originalURL)
			if err == nil {

				return existing, nil
			}
			return "", wrapError(op, err)
		}

		return "", wrapError(op, err)
	}

	return "", wrapError(op, ErrAttemptsExceeded)
}
