package evm_wallet

import (
	"context"

	evmWalletModel "github.com/zd4r/wallet-manager/internal/model/evm_wallet"
)

type evmWalletStore interface {
	Create(ctx context.Context, wallet *evmWalletModel.EvmWallet) (int64, error)
	GetByID(ctx context.Context, id int64) (*evmWalletModel.EvmWallet, error)
	GetList(ctx context.Context) ([]evmWalletModel.EvmWallet, error)
	DeleteByID(ctx context.Context, id int64) error
}

type passphraseStore interface {
	Set(val []byte)
	Get() string
}
