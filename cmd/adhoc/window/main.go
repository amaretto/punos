package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	fmt.Printf("out: %#v\n", string(out))
	result := string(out)
	result = strings.TrimRight(result, "\n")
	fmt.Println(result)
	hoge := strings.Split(result, " ")
	fmt.Println(hoge)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		log.Fatal(err)
	}
}
