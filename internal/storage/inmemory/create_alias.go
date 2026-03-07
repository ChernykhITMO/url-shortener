package inmemory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) Create(_ context.Context, alias, originalURL string) (string, error) {
	const op = "inmemory.Create"

	s.mux.Lock()
	defer s.mux.Unlock()

	if existingAlias, ok := s.urlToAlias[originalURL]; ok {
		return existingAlias, nil
	}

	if existingURL, ok := s.aliasToURL[alias]; ok && existingURL != originalURL {
		return "", fmt.Errorf("%s: %w", op, storage.ErrAliasConflict)
	}

	s.urlToAlias[originalURL] = alias
	s.aliasToURL[alias] = originalURL

	return alias, nil
}
