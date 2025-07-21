package main

import "fmt"

func main() {
	var n int64
	// Read the input n
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// Hexagonal number formula: h_n = 2n^2 - n
	h := 2*n*n - n
	fmt.Println(h)
}
