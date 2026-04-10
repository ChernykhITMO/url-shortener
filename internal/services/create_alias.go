package services

import (
	"context"
	"errors"
	"strings"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

var ErrAttemptsExceeded = errors.New("attempts exceeded")

func (s *Service) CreateAlias(ctx context.Context, originalURL, requestedAlias string) (string, error) {
	const op = "services.CreateAlias"

	normalizedURL, err := normalizeAndValidateURL(originalURL)
	if err != nil {
		return "", wrapError(op, err)
	}

	requestedAlias = strings.TrimSpace(requestedAlias)
	if requestedAlias != "" {
		if err := validateAlias(requestedAlias); err != nil {
			return "", wrapError(op, err)
		}

		createdAlias, err := s.storage.Create(ctx, requestedAlias, normalizedURL)
		if err == nil {
			return createdAlias, nil
		}
		if errors.Is(err, storage.ErrAliasConflict) {
			return "", wrapError(op, ErrAliasTaken)
		}
		if errors.Is(err, storage.ErrInvalidAlias) {
			return "", wrapError(op, ErrInvalidAlias)
		}
		return "", wrapError(op, err)
	}

	for i := 0; i < s.maxAttempts; i++ {
		alias, err := s.generateAlias()
		if err != nil {
			return "", wrapError(op, err)
		}

		createdAlias, err := s.storage.Create(ctx, alias, normalizedURL)
		if err == nil {
			return createdAlias, nil
		}

		if errors.Is(err, storage.ErrAliasConflict) {
			continue
		}

		if errors.Is(err, storage.ErrInvalidAlias) {
			return "", wrapError(op, ErrInvalidAlias)
		}

		return "", wrapError(op, err)
	}

	return "", wrapError(op, ErrAttemptsExceeded)
}
