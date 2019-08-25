package main

import (
	"log"

	"github.com/amaretto/punos/ui"
)

func doUI(logger *log.Logger) error {
	app := ui.NewApp(client, url)
	app.SetLogger(logger)

	app.Run()
	return nil
}
