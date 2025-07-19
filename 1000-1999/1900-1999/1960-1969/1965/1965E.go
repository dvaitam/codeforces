package main

import (
	"bufio"
	"fmt"
	"os"
)

// Move represents a single operation with coordinates and color.
type Move struct {
	x, y, z, c int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	// read matrix a, 1-indexed
	a := make([][]int, n+2)
	for i := range a {
		a[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	var ans []Move
	// first part: process each cell
	for i := n; i >= 1; i-- {
		for j := 1; j <= m; j++ {
			// l from 2 to n-i+1
			for l := 2; l <= n-i+1; l++ {
				ans = append(ans, Move{i, j, l, a[i][j]})
			}
			// l from i+1 to 2*i-1
			for l := i + 1; l <= 2*i-1; l++ {
				ans = append(ans, Move{l, j, n - i + 1, a[i][j]})
			}
			// l from n-i+2 to n + a[i][j]
			for l := n - i + 2; l <= n+a[i][j]; l++ {
				ans = append(ans, Move{2*i - 1, j, l, a[i][j]})
			}
		}
	}
	// second part: for each t
	for t := 1; t <= k; t++ {
		for i := 1; i < n; i++ {
			if (i&1) == 1 || i+1 == n {
				for j := 1; j <= m; j++ {
					ans = append(ans, Move{2 * i, j, t + n - 1, t})
				}
			}
		}
		for i := 1; i < n; i++ {
			ans = append(ans, Move{2 * i, m + 1, t + n - 1, t})
			if i+1 < n {
				ans = append(ans, Move{2*i + 1, m + 1, t + n - 1, t})
			}
		}
	}
	// output
	fmt.Fprintln(writer, len(ans))
	for _, mv := range ans {
		fmt.Fprintln(writer, mv.x, mv.y, mv.z, mv.c)
	}
}
