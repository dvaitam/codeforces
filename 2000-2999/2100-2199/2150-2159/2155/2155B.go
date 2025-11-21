package main

import (
	"bufio"
	"fmt"
	"os"
)

type cell struct {
	r, c int
}

func direction(from, to cell) byte {
	if from.r == to.r {
		if to.c == from.c+1 {
			return 'R'
		}
		if to.c == from.c-1 {
			return 'L'
		}
	}
	if from.c == to.c {
		if to.r == from.r+1 {
			return 'D'
		}
		if to.r == from.r-1 {
			return 'U'
		}
	}
	panic("non-adjacent cells passed to direction")
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		if k == n*n-1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		fmt.Fprintln(out, "YES")
		grid := make([][]byte, n)
		for i := range grid {
			grid[i] = make([]byte, n)
		}

		order := make([]cell, 0, n*n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				for j := 0; j < n; j++ {
					order = append(order, cell{i, j})
				}
			} else {
				for j := n - 1; j >= 0; j-- {
					order = append(order, cell{i, j})
				}
			}
		}

		for idx := 0; idx < k; idx++ {
			c := order[idx]
			grid[c.r][c.c] = 'U'
		}

		region := order[k:]
		if len(region) > 0 {
			if len(region) == 1 {
				panic("unexpected region of size 1")
			}
			grid[region[0].r][region[0].c] = direction(region[0], region[1])
			grid[region[1].r][region[1].c] = direction(region[1], region[0])
			for idx := 2; idx < len(region); idx++ {
				a := region[idx]
				b := region[idx-1]
				grid[a.r][a.c] = direction(a, b)
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprintln(out, string(grid[i]))
		}
	}
}
