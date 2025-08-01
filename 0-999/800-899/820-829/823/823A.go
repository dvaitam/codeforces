package main

import (
	"fmt"
)

const MAXN = 300001

func main() {
	var n, k int
	fmt.Scanf("%d %d", &n, &k)

	// Initialize array a
	a := make([]int, MAXN)

	// Calculate d
	d := n - 1 - k

	for i := 0; i < k; i++ {
		a[i] = (d + i)/k + 1
	}

	fmt.Println(a[k-2] + a[k-1])

	// Print edges
	c := 1
	for i := 0; i < k; i++ {
		for j := 0; j+1 < a[i]; j++ {
			fmt.Printf("%d %d\n", c, c+1)
			c++
		}
		fmt.Printf("%d %d\n", c, n)
		c++
	}
}
