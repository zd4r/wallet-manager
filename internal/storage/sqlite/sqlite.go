package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

const (
	op         = "storage.sqlite.New"
	driverName = "sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewWithContext(ctx context.Context, storagePath string) (*Storage, error) {
	db, err := sql.Open(driverName, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.SetDialect(driverName); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.RunContext(ctx, "up", db, "."); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) Check() (interface{}, error) {
	return s.db.Stats(), s.db.Ping()
}

func (s *Storage) Prepare(query string) (*sql.Stmt, error) {
	return s.db.Prepare(query)
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
