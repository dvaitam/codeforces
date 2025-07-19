package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// g[i][j] stores the j-th capacity of i-th row (1-based)
	var g [5][5]int
	for i := 1; i <= 4; i++ {
		fmt.Scan(&g[i][1], &g[i][2], &g[i][3], &g[i][4])
	}
	found := false
	pos, a, b := 0, 0, 0
	for i := 1; i <= 4; i++ {
		// compute mins
		m1 := g[i][1]
		if g[i][2] < m1 {
			m1 = g[i][2]
		}
		m2 := g[i][3]
		if g[i][4] < m2 {
			m2 = g[i][4]
		}
		if m1+m2 <= n {
			found = true
			pos = i
			a = m1
			b = n - m1
			break
		}
	}
	if found {
		fmt.Println(pos, a, b)
	} else {
		fmt.Println(-1)
	}
}
