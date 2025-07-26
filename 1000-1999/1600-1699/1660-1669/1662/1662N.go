package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	type cell struct{ r, c int }
	ans := int64(0)
	for r1 := 0; r1 < n; r1++ {
		for r2 := r1 + 1; r2 < n; r2++ {
			for c1 := 0; c1 < n; c1++ {
				for c2 := c1 + 1; c2 < n; c2++ {
					vals := []struct {
						v int
						p cell
					}{
						{grid[r1][c1], cell{r1, c1}},
						{grid[r1][c2], cell{r1, c2}},
						{grid[r2][c1], cell{r2, c1}},
						{grid[r2][c2], cell{r2, c2}},
					}
					sort.Slice(vals, func(i, j int) bool { return vals[i].v < vals[j].v })
					p1 := vals[0].p
					p2 := vals[1].p
					if p1.r == p2.r || p1.c == p2.c {
						ans++
					}
				}
			}
		}
	}
	fmt.Println(ans)
}
