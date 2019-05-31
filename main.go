package main

import "fmt"

func main() {
	fmt.Println("\n\n")
	fmt.Println("_ __  _   _ _ __   ___  ___")
	fmt.Println("| '_ \\| | | | '_ \\ / _ \\/ __|")
	fmt.Println("| |_) | |_| | | | | (_) \\__ \\")
	fmt.Println("| .__/ \\__,_|_| |_|\\___/|___/")
	fmt.Println("|_|")
	fmt.Println("\n\n")

	var in string
	fmt.Printf("[empty]>>")
	fmt.Scanln(&in)

	fmt.Println(in)

}
