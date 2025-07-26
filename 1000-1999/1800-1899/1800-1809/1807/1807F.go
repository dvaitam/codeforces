package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct{ r, c, dx, dy int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, i1, j1, i2, j2 int
		var d string
		fmt.Fscan(in, &n, &m, &i1, &j1, &i2, &j2, &d)

		dir := map[string][2]int{"DR": {1, 1}, "DL": {1, -1}, "UR": {-1, 1}, "UL": {-1, -1}}
		dx, dy := dir[d][0], dir[d][1]
		r, c := i1, j1
		bounces := 0
		visited := make(map[state]bool)

		for {
			if r == i2 && c == j2 {
				fmt.Fprintln(out, bounces)
				break
			}
			st := state{r, c, dx, dy}
			if visited[st] {
				fmt.Fprintln(out, -1)
				break
			}
			visited[st] = true

			nr, nc := r+dx, c+dy
			bounce := false
			if nr < 1 || nr > n {
				dx = -dx
				bounce = true
			}
			if nc < 1 || nc > m {
				dy = -dy
				bounce = true
			}
			if bounce {
				bounces++
			}
			r += dx
			c += dy
		}
	}
}
