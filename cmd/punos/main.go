package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/amaretto/punos/ui"
)

var version = "1.0.0"

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()
	if showVersion {
		fmt.Println("version:", version)
		return
	}

	var dlog *log.Logger
	logfile := "hoge"
	if logfile != "" {
		f, e := os.Create(logfile)
		if e == nil {
			dlog = log.New(f, "DEBUG:", log.LstdFlags)
			log.SetOutput(f)
		}
	}

	if e := doUI(dlog); e != nil {
		os.Exit(1)
	}
}

func doUI(logger *log.Logger) error {

	app := ui.NewApp()
	app.SetLogger(logger)

	app.Run()
	return nil
}
