package main

import (
	"fmt"
)

func expected(a, b int) int {
	n := a + b
	count := 0
	for k := 1; k <= n; k++ {
		m := n / k
		if m == 0 {
			continue
		}
		r := n % k
		for x := 0; x <= r; x++ {
			rem := a - x*(m+1)
			if rem < 0 {
				continue
			}
			if rem%m != 0 {
				continue
			}
			y := rem / m
			if y >= 0 && y <= k-r {
				count++
				break
			}
		}
	}
	return count
}

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	fmt.Println(expected(a, b))
}
