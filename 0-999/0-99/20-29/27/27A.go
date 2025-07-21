package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// mark presence of indices up to n+1
	present := make([]bool, n+2)
	for i := 0; i < n; i++ {
		var x int
		if _, err := fmt.Scan(&x); err != nil {
			return
		}
		if x >= 1 && x <= n+1 {
			present[x] = true
		}
	}
	// find smallest missing positive integer
	for i := 1; i <= n+1; i++ {
		if !present[i] {
			fmt.Println(i)
			return
		}
	}
}
