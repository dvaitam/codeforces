package main

import (
	"fmt"
	"math"
)

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// Compute maximum k such that k*(k+1)/2 <= n
	maxk := int((math.Sqrt(1+8*float64(n)) - 1) / 2)
	// Generate triangular numbers T_i = i*(i+1)/2 for i = 1..maxk
	t := make([]int64, maxk)
	for i := 1; i <= maxk; i++ {
		t[i-1] = int64(i*(i+1)) / 2
	}
	// Two-pointer to find if any pair sums to n
	l, r := 0, len(t)-1
	for l <= r {
		s := t[l] + t[r]
		if s == n {
			fmt.Println("YES")
			return
		}
		if s < n {
			l++
		} else {
			r--
		}
	}
	fmt.Println("NO")
}
