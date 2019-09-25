package ui

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetWindowSize return window height and width utilizing stty
func GetWindowSize() (height, width int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	result := string(out)
	result = strings.TrimRight(result, "\n")
	window := strings.Split(result, " ")
	if err != nil {
		log.Fatal(err)
	}

	height, _ = strconv.Atoi(window[0])
	width, _ = strconv.Atoi(window[1])

	return
}
