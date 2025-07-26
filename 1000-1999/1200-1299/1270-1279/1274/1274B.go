package main

import "fmt"

func main() {
	for i := 0; i < 37; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(1)
	}
}
