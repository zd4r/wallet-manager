package evm_wallet

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmWalletModel "github.com/zd4r/wallet-manager/internal/model/evm_wallet"
)

type Service struct {
	walletStore     evmWalletStore
	keyStore        *keystore.KeyStore
	passphraseStore passphraseStore
}

func New(walletStore evmWalletStore, keyStore *keystore.KeyStore, passphraseStore passphraseStore) *Service {
	return &Service{
		walletStore:     walletStore,
		keyStore:        keyStore,
		passphraseStore: passphraseStore,
	}
}

func (s *Service) Create(ctx context.Context, wallet *evmWalletModel.EvmWallet) (int64, error) {
	pk, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return 0, fmt.Errorf("failed to parse private key: %w", err)
	}

	// TODO: make transaction with db or smth else, now there is a risk of desync
	account, err := s.keyStore.ImportECDSA(pk, s.passphraseStore.Get())
	if err != nil && err != keystore.ErrAccountAlreadyExists {
		return 0, fmt.Errorf("failed to save wallet to key store: %w", err)
	}

	wallet.Address = account.Address.String()

	return s.walletStore.Create(ctx, wallet)
}

func (s *Service) GetList(ctx context.Context) ([]evmWalletModel.EvmWallet, error) {
	return s.walletStore.GetList(ctx)
}

func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	wallet, err := s.walletStore.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get wallet [id=%d] by id: %w", id, err)
	}

	account, err := s.keyStore.Find(accounts.Account{Address: common.HexToAddress(wallet.Address)})
	if err != nil {
		return fmt.Errorf(
			"failed to find wallet [address=%s] in keystore: %w",
			wallet.GetShortAddress(),
			err,
		)
	}

	// TODO: make transactional with db or sync db on app start up, now there is a risk of desync
	if err := s.keyStore.Delete(account, s.passphraseStore.Get()); err != nil {
		return fmt.Errorf(
			"failed to delete account [address=%s] from keystore: %w",
			account.Address.String(),
			err,
		)
	}

	return s.walletStore.DeleteByID(ctx, id)
}

func (s *Service) CheckAccess() error {
	if err := s.keyStore.TimedUnlock(
		s.keyStore.Accounts()[rand.Intn(len(s.keyStore.Accounts()))],
		s.passphraseStore.Get(),
		1*time.Microsecond,
	); err != nil {
		return fmt.Errorf("failed to unlock keystore: %w", err)
	}

	return nil
}
