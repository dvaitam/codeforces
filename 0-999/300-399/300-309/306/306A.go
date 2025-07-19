package main

import (
	"fmt"
)

func main() {
	var n, m int
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}
	each := n / m
	mod := n % m
	rem := m - mod
	for i := 0; i < mod; i++ {
		fmt.Printf("%d ", each+1)
	}
	for i := 0; i < rem; i++ {
		fmt.Printf("%d ", each)
	}
}
