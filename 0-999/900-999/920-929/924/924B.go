package main

import (
	"fmt"
)

func main() {
	var n int
	var U int
	if _, err := fmt.Scan(&n, &U); err != nil {
		return
	}
	E := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&E[i])
	}
	best := -1.0
	k := 0
	for i := 0; i < n; i++ {
		if k < i {
			k = i
		}
		for k+1 < n && E[k+1]-E[i] <= U {
			k++
		}
		if k >= i+2 {
			val := float64(E[k]-E[i+1]) / float64(E[k]-E[i])
			if val > best {
				best = val
			}
		}
	}
	if best < 0 {
		fmt.Println(-1)
	} else {
		fmt.Printf("%.12f\n", best)
	}
}
