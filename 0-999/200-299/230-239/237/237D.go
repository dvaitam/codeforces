package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	// Read n-1 edges
	edges := make([][2]int, n)
	c := make([]int, n+1)
	for i := 1; i < n; i++ {
		var x, y int
		fmt.Scan(&x, &y)
		edges[i][0], edges[i][1] = x, y
		if c[x] == 0 {
			c[x] = i
		}
		if c[y] == 0 {
			c[y] = i
		}
	}

	// Output
	fmt.Println(n)
	for i := 1; i < n; i++ {
		fmt.Printf("2 %d %d\n", edges[i][0], edges[i][1])
	}
	for i := 1; i < n; i++ {
		u, v := edges[i][0], edges[i][1]
		if i != c[u] {
			fmt.Printf("%d %d\n", i, c[u])
		}
		if i != c[v] {
			fmt.Printf("%d %d\n", i, c[v])
		}
	}
}
