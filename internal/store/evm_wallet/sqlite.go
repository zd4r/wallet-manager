package evm_wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	evmWalletModel "github.com/zd4r/wallet-manager/internal/model/evm_wallet"
	"github.com/zd4r/wallet-manager/internal/storage/sqlite"
	"github.com/zd4r/wallet-manager/internal/store"
)

type SQLite struct {
	db *sqlite.Storage
}

func New(db *sqlite.Storage) *SQLite {
	return &SQLite{db}
}

// Create saves wallet to db
func (s *SQLite) Create(ctx context.Context, wallet *evmWalletModel.EvmWallet) (int64, error) {
	const op = "store.wallet.sqlite.Create"

	stmt, err := s.db.Prepare("INSERT INTO evm_wallet(name, address) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, wallet.Name, wallet.Address)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, store.EntityAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *SQLite) GetByID(ctx context.Context, id int64) (*evmWalletModel.EvmWallet, error) {
	const op = "store.wallet.sqlite.GetByID"

	wallet := new(evmWalletModel.EvmWallet)

	stmt, err := s.db.Prepare("SELECT id, name, address FROM evm_wallet WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&wallet.ID, &wallet.Name, &wallet.Address); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, store.EntityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return wallet, nil

}

func (s *SQLite) GetList(ctx context.Context) ([]evmWalletModel.EvmWallet, error) {
	const op = "store.wallet.sqlite.GetList"

	stmt, err := s.db.Prepare("SELECT id, name, address FROM evm_wallet")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var wallets []evmWalletModel.EvmWallet
	for rows.Next() {
		var wallet evmWalletModel.EvmWallet

		if err := rows.Scan(&wallet.ID, &wallet.Name, &wallet.Address); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (s *SQLite) DeleteByID(ctx context.Context, id int64) error {
	const op = "store.wallet.sqlite.DeleteById"

	stmt, err := s.db.Prepare("DELETE FROM evm_wallet WHERE id=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
