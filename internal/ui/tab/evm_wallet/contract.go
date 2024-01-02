package evm_wallet

import (
	"context"

	evmWalletModel "github.com/zd4r/wallet-manager/internal/model/evm_wallet"
)

type walletService interface {
	Create(ctx context.Context, wallet *evmWalletModel.EvmWallet) (int64, error)
	GetList(ctx context.Context) ([]evmWalletModel.EvmWallet, error)
	DeleteByID(ctx context.Context, id int64) error
}
