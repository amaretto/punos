package main

import (
	"fmt"
	"math"
)

func main() {
	var omega, alpha float64
	var q float64
	q = 4

	omega = 2.0 * 3.14159265
	alpha = math.Sin(omega) / (2.0 * q)
	fmt.Println(alpha)
	alpha = math.Sin(omega) / (float64(2.0) * q)
	fmt.Println(alpha)
	//	fmt.Println(omega)
	//
	//	omega = float64(2.0) * float64(3.14159265)
	//	fmt.Println(omega)

}
