package main

import (
	"log"
	"os"
	"runtime/trace"

	"github.com/amaretto/punos/ui"
)

func main() {

	var dlog *log.Logger
	logfile := "hoge"
	if logfile != "" {
		f, e := os.Create(logfile)
		if e == nil {
			dlog = log.New(f, "DEBUG:", log.LstdFlags)
			log.SetOutput(f)
		}
	}

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()
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
