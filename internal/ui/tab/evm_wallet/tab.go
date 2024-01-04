package evm_wallet

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"
	evmWalletModel "github.com/zd4r/wallet-manager/internal/model/evm_wallet"
)

var columns = []fyne.CanvasObject{
	widget.NewLabel("id"),
	widget.NewLabel("name"),
	widget.NewLabel("address"),
	widget.NewLabel("actions"),
}

type Tab struct {
	mainWindow    fyne.Window
	walletService walletService
}

func New(mainWindow fyne.Window, walletService walletService) *Tab {
	return &Tab{
		mainWindow:    mainWindow,
		walletService: walletService,
	}
}

func (t *Tab) Build(ctx context.Context) (*fyne.Container, error) {
	// wallet list
	walletList, err := t.walletService.GetList(ctx)
	if err != nil {
		return nil, err
	}

	// wallet list widget
	walletListWidget := widget.NewList(
		func() int {
			return len(walletList)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(len(columns)),
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewToolbar(
					widget.NewToolbarAction(theme.DeleteIcon(), func() {}),
					widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
				),
			)
		},
		nil,
	)

	buildWalletToolbar := func(id int) *widget.Toolbar {
		return widget.NewToolbar(
			widget.NewToolbarAction(theme.DeleteIcon(), func() {
				if len(walletList) <= id {
					return
				}

				if err = t.walletService.DeleteByID(ctx, walletList[id].ID); err != nil {
					dialog.ShowInformation(
						"error occurred",
						err.Error(),
						t.mainWindow,
					)
					return
				}

				walletList, err = t.walletService.GetList(ctx)
				if err != nil {
					dialog.ShowInformation(
						"error occurred",
						err.Error(),
						t.mainWindow,
					)
					return
				}

				walletListWidget.Refresh()
			}),
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
				t.mainWindow.Clipboard().SetContent(walletList[id].Address)
			}),
			widget.NewToolbarAction(theme.MailAttachmentIcon(), func() {
				if err := browser.OpenURL(
					fmt.Sprintf("https://etherscan.io/address/%s", walletList[id].Address),
				); err != nil {
					dialog.ShowInformation(
						"error occurred",
						err.Error(),
						t.mainWindow,
					)
				}
			}),
		)
	}

	walletListWidget.UpdateItem = func(id widget.ListItemID, object fyne.CanvasObject) {
		c := object.(*fyne.Container)

		idLabel := c.Objects[0].(*widget.Label)
		idLabel.SetText(fmt.Sprintf("%d", walletList[id].ID))

		nameLabel := c.Objects[1].(*widget.Label)
		nameLabel.SetText(walletList[id].Name)

		addressLabel := c.Objects[2].(*widget.Label)
		addressLabel.SetText(walletList[id].GetShortAddress())

		c.Objects[3] = buildWalletToolbar(id)
	}

	walletListWidget.OnSelected = func(id widget.ListItemID) {
		walletListWidget.Unselect(id)
	}

	// wallet input widget
	walletNameInput := widget.NewEntry()
	walletNameInput.SetPlaceHolder("wallet name")

	pkInput := widget.NewPasswordEntry()
	pkInput.SetPlaceHolder("private key")

	walletInput := container.New(
		layout.NewGridLayout(2),
		walletNameInput,
		pkInput,
	)

	inputWidget := container.NewBorder(nil, nil, nil,
		widget.NewToolbar(
			widget.NewToolbarAction(
				theme.DocumentSaveIcon(),
				func() {
					defer func() {
						pkInput.Text = ""
						pkInput.Refresh()
					}()
					defer walletListWidget.Refresh()

					wallet := &evmWalletModel.EvmWallet{
						Name:       walletNameInput.Text,
						PrivateKey: pkInput.Text,
					}
					if err := wallet.Validate(); err != nil {
						dialog.ShowInformation(
							"error occurred",
							err.Error(),
							t.mainWindow,
						)
						return
					}

					_, err := t.walletService.Create(ctx, wallet)
					if err != nil {
						dialog.ShowInformation(
							"error occurred",
							err.Error(),
							t.mainWindow,
						)
						return
					}

					walletList, err = t.walletService.GetList(ctx)
					if err != nil {
						dialog.ShowInformation(
							"error occurred",
							err.Error(),
							t.mainWindow,
						)
						return
					}
				},
			),
		),
		walletInput,
	)

	// result tab
	c := container.NewBorder(
		container.NewVBox(widget.NewSeparator(), inputWidget),
		nil, nil, nil,
		walletListWidget,
	)

	return c, nil
}
