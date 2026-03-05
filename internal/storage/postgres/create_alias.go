package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation             = "23505"
	constraintAliasPK           = "url_pkey"
	constraintOriginalURLUnique = "url_original_url_key"
)

func (s *Storage) Create(ctx context.Context, alias, originalURL string) error {
	const op = "postgres.Create"

	const query = `INSERT INTO url(alias, original_url) VALUES ($1, $2)`

	if _, err := s.db.ExecContext(ctx, query, alias, originalURL); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			switch pgErr.ConstraintName {
			case constraintAliasPK:
				return fmt.Errorf("%s: %w", op, storage.ErrAliasConflict)
			case constraintOriginalURLUnique:
				return fmt.Errorf("%s: %w", op, storage.ErrURLConflict)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
