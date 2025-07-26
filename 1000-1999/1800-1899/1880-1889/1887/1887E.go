package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type cell struct{ x, y int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		cells := make([]cell, 2*n+1)
		rowMap := make(map[int]map[int]int)
		for i := 1; i <= 2*n; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			cells[i] = cell{x, y}
			if rowMap[x] == nil {
				rowMap[x] = make(map[int]int)
			}
			rowMap[x][y] = i
		}
		rows := make([]int, 0, len(rowMap))
		for r := range rowMap {
			rows = append(rows, r)
		}
		sort.Ints(rows)
		found := false
		var ans [4]cell
		for i := 0; i < len(rows) && !found; i++ {
			for j := i + 1; j < len(rows) && !found; j++ {
				r1, r2 := rows[i], rows[j]
				cols := make([]struct{ c, i1, i2 int }, 0)
				for c1, idx1 := range rowMap[r1] {
					if idx2, ok := rowMap[r2][c1]; ok {
						cols = append(cols, struct{ c, i1, i2 int }{c1, idx1, idx2})
					}
				}
				if len(cols) < 2 {
					continue
				}
				for a := 0; a < len(cols) && !found; a++ {
					for b := a + 1; b < len(cols) && !found; b++ {
						i1 := cols[a].i1
						i2 := cols[b].i1
						i3 := cols[a].i2
						i4 := cols[b].i2
						if i1 == i2 || i1 == i3 || i1 == i4 || i2 == i3 || i2 == i4 || i3 == i4 {
							continue
						}
						ans = [4]cell{cells[i1], cells[i2], cells[i3], cells[i4]}
						found = true
					}
				}
			}
		}
		if found {
			fmt.Fprintln(out, "Yes")
			for _, c := range ans {
				fmt.Fprintf(out, "%d %d\n", c.x, c.y)
			}
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
