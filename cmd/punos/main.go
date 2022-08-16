package main

import (
	"os"

	"github.com/amaretto/punos/pkg/cmd/punos"
)

func main() {
	// default command is player
	if len(os.Args) == 1 {
		os.Args = append([]string{os.Args[0], "player"}, os.Args[1:]...)
	}
	punos.NewCommand().Execute()
}
