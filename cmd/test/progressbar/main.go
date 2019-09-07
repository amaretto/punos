package main

import "fmt"

func main() {
	var width, position, length int
	width = 100
	position = 33
	length = 100
	fmt.Println(getProgressbar(width, position, length))
}

func getProgressbar(width, position, length int) string {
	current := float64(width) * float64(position) / float64(length)
	str := "["
	for i := 0; i < width; i++ {
		if i > int(current) {
			str = str + "-"
		} else {
			str = str + "="
		}
	}
	str = str + "]"
	return str
}
