package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/zd4r/wallet-manager/migrations"

	"github.com/zd4r/wallet-manager/internal/app"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	a := app.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
