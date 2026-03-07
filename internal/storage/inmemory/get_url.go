package in_memory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "in_memory.GetURL"

	s.mux.RLock()
	url, ok := s.aliasToURL[alias]
	s.mux.RUnlock()

	if ok {
		return url, nil
	}

	return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
}
