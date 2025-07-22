package main

import (
	"fmt"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var l, r, a int
	if _, err := fmt.Scan(&l, &r, &a); err != nil {
		return
	}
	total := l + r + a
	half := total / 2
	// each side cannot exceed its available count including ambidexters
	half = min(half, l+a)
	half = min(half, r+a)
	fmt.Println(2 * half)
}
