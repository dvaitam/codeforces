package main

import (
	"fmt"
)

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	fmt.Println(a * b)
}
