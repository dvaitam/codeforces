package main

import (
	"bufio"
	"fmt"
	"os"
)

// The interactive version of the problem asks to reveal edges
// using queries.  In this archive we are provided the whole
// graph in the input, so we simply build its connected
// components and assign each component a unique color.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	g := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	color := make([]int, n+1)
	cur := 0
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if color[i] != 0 {
			continue
		}
		cur++
		q = q[:0]
		q = append(q, i)
		color[i] = cur
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for _, v := range g[u] {
				if color[v] == 0 {
					color[v] = cur
					q = append(q, v)
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, color[i])
	}
	out.WriteByte('\n')
}
