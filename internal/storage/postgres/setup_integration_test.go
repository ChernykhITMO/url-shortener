//go:build integration
// +build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ChernykhITMO/url-shortener/internal/config"
	"github.com/pressly/goose/v3"
)

var (
	testStore *Storage
	testDB    *sql.DB
)

func TestMain(m *testing.M) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s, err := New(ctx, config.Postgres{
		DSN:             dsn,
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		panic(err)
	}
	testStore = s

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	testDB = db

	if err := runGooseUp(db); err != nil {
		panic(err)
	}

	code := m.Run()
	_ = db.Close()
	_ = s.Close()
	os.Exit(code)
}

func cleanupURLTable(t *testing.T) {
	t.Helper()

	if _, err := testDB.Exec(`TRUNCATE TABLE url`); err != nil {
		t.Fatalf("truncate url failed: %v", err)
	}
}

func runGooseUp(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("cannot resolve current file path")
	}

	migrationsDir := filepath.Join(filepath.Dir(file), "../../../migrations/migrations")
	return goose.Up(db, migrationsDir)
}

