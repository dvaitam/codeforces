package main

import (
	"fmt"
)

func main() {
	var a, b int
	for {
		if _, err := fmt.Scan(&a, &b); err != nil {
			return
		}
		fmt.Println(a + b)
	}
}
