package main

import "fmt"

func main() {
	var t int
	fmt.Scan(&t)
	for ; t > 0; t-- {
		var x, y int
		fmt.Scan(&x, &y)
		fmt.Println(x / y) // intentional runtime error when y == 0
	}
}
