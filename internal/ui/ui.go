package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	evmWalletTab "github.com/zd4r/wallet-manager/internal/ui/tab/evm_wallet"
)

const (
	appName = "wallet-manager"
)

func New(ctx context.Context, walletService evmWalletService) (fyne.App, error) {
	// init ui
	a := app.New()
	mainWindow := a.NewWindow(appName)
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.Size{Width: 750, Height: 500})

	// evm wallet tab
	ewt := evmWalletTab.New(mainWindow, walletService)
	ewtContent, err := ewt.Build(ctx)
	if err != nil {
		return nil, err
	}

	// finalise
	tabs := container.NewAppTabs(
		container.NewTabItem("evm wallets", ewtContent),
	)

	mainWindow.SetContent(tabs)
	mainWindow.Show()

	return a, nil
}
