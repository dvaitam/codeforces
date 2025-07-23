package main

import (
	"fmt"
)

func comb5(n int64) int64 {
	return n * (n - 1) * (n - 2) * (n - 3) * (n - 4) / 120
}

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	c := comb5(n)
	ans := c * c * 120
	fmt.Println(ans)
}
