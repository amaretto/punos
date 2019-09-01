package main

import (
	"log"

	"github.com/amaretto/punos/ui"
)

func doUI(logger *log.Logger) error {

	app := ui.NewApp()
	app.SetLogger(logger)

	app.Run()
	return nil
}
