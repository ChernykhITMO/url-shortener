package in_memory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetAliasByURL(ctx context.Context, originalURL string) (string, error) {
	const op = "in_memory.GetAliasByURL"

	s.mux.RLock()
	alias, ok := s.urlToAlias[originalURL]
	s.mux.RUnlock()

	if !ok {
		return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return alias, nil
}
