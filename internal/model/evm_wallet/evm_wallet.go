package evm_wallet

import (
	"errors"
	"fmt"
)

type EvmWallet struct {
	ID         int64
	Name       string
	Address    string
	PrivateKey string
}

func (w *EvmWallet) GetShortAddress() string {
	l := len(w.Address)
	return fmt.Sprintf("%s...%s", w.Address[:6], w.Address[l-4:l])
}

func (w *EvmWallet) Validate() error {
	if w.Name == "" {
		return errors.New("invalid wallet name: empty")
	}

	return nil
}
