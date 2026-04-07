package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetAlias(ctx context.Context, originalURL string) (string, error) {
	const op = "postgres.GetAlias"

	const query = `SELECT alias FROM url WHERE original_url = $1`

	var alias string
	if err := s.db.QueryRowContext(ctx, query, originalURL).Scan(&alias); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return alias, nil
}
