package main

import (
	"fmt"
)

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	if a == 0 {
		var s []int
		fmt.Println(s[0])
		return
	}
	fmt.Println(a + b)
}
