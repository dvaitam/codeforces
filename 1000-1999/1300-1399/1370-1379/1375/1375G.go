package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	q := []int{0}
	color[0] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if color[v] == -1 {
				color[v] = color[u] ^ 1
				q = append(q, v)
			}
		}
	}

	c0, c1 := 0, 0
	for _, c := range color {
		if c == 0 {
			c0++
		} else {
			c1++
		}
	}
	if c1 < c0 {
		c0 = c1
	}
	if c0 > 0 {
		c0--
	}
	fmt.Fprintln(out, c0)
}
