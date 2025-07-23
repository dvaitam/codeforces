package main

import (
	"fmt"
)

func main() {
	var a uint
	if _, err := fmt.Scan(&a); err != nil {
		return
	}
	fmt.Println(uint64(1) << a)
}
