package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	evmWalletService "github.com/zd4r/wallet-manager/internal/service/evm_wallet"
	"github.com/zd4r/wallet-manager/internal/storage/sqlite"
	evmWalletStore "github.com/zd4r/wallet-manager/internal/store/evm_wallet"
	"github.com/zd4r/wallet-manager/internal/store/passphrase"
	"github.com/zd4r/wallet-manager/internal/ui"
	"golang.org/x/term"
)

type App struct{}

func New() *App {
	return &App{}
}

const (
	keystoreDirPath = "/.ethereum/keystore"
	databasePath    = "/.ethereum/wallet-manager.db"
)

func (a *App) Run() error {
	cxt := context.Background()

	// init global passphrase store
	pp := passphrase.New()

	// set passphrase
	password, err := readPassphrase()
	if err != nil {
		return err
	}
	pp.Set(password)

	// get $HOME
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create .ethereum dir
	path := fmt.Sprintf("%s/.ethereum", homeDir)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// init keystore
	ks := keystore.NewKeyStore(
		fmt.Sprintf("%s%s", homeDir, keystoreDirPath), // TODO: change path to $HOME/.ethereum/keystore
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	// init bot data storage
	storage, err := sqlite.NewWithContext(cxt, fmt.Sprintf("%s%s", homeDir, databasePath))
	if err != nil {
		return err
	}
	defer storage.Stop()

	// init wallet service
	evmWalletSrv := evmWalletService.New(
		evmWalletStore.New(storage),
		ks,
		pp,
	)
	if err := evmWalletSrv.CheckAccess(); err != nil {
		return err
	}

	// init ui
	appUI, err := ui.New(cxt, evmWalletSrv)
	if err != nil {
		return err
	}

	// run app
	appUI.Run()

	return nil
}

func readPassphrase() ([]byte, error) {
	defer fmt.Println()

	fmt.Print("password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to term.ReadPassword: %w", err)
	}

	return password, nil
}
