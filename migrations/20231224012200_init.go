package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS evm_wallet (
				    id INTEGER PRIMARY KEY,
				    name TEXT NOT NULL UNIQUE,
				    address TEXT NOT NULL UNIQUE
				);`,
	)

	return err
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		`DROP TABLE IF EXISTS evm_wallet;`,
	)

	return err
}
