package main

import (
	"fmt"
)

func uniqueDigits(x int) bool {
	seen := [10]bool{}
	if x == 0 {
		seen[0] = true
		return true
	}
	for x > 0 {
		d := x % 10
		if seen[d] {
			return false
		}
		seen[d] = true
		x /= 10
	}
	return true
}

func main() {
	var l, r int
	if _, err := fmt.Scan(&l, &r); err != nil {
		return
	}
	for x := l; x <= r; x++ {
		if uniqueDigits(x) {
			fmt.Println(x)
			return
		}
	}
	fmt.Println(-1)
}
