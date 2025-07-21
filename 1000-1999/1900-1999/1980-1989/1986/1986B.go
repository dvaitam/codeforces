package main

import (
	"bufio"
	"fmt"
	"os"
)

type cell struct{ r, c int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &a[i][j])
			}
		}

		q := make([]cell, 0)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				q = append(q, cell{i, j})
			}
		}
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			r, c := cur.r, cur.c
			mx := -1 << 60
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < n && nc >= 0 && nc < m {
					if a[nr][nc] > mx {
						mx = a[nr][nc]
					}
				}
			}
			if mx == -1<<60 { // isolated single cell
				continue
			}
			if a[r][c] > mx {
				a[r][c] = mx
				q = append(q, cell{r, c})
				for _, d := range dirs {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < n && nc >= 0 && nc < m {
						q = append(q, cell{nr, nc})
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, a[i][j])
			}
			writer.WriteByte('\n')
		}
	}
}
