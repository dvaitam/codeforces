package main

import (
	"fmt"
)

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	if a > b {
		fmt.Println(a)
	} else {
		fmt.Println(b)
	}
}
