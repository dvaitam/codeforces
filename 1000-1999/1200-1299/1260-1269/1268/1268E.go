package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

type State struct {
	v    int
	last int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g := make([][]Edge, n+1)
	for i := 1; i <= m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		g[a] = append(g[a], Edge{to: b, w: i})
		g[b] = append(g[b], Edge{to: a, w: i})
	}

	ans := make([]int, n+1)
	for start := 1; start <= n; start++ {
		visited := make([]bool, n+1)
		q := make([]State, 0)
		q = append(q, State{v: start, last: 0})
		visited[start] = true
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			for _, e := range g[cur.v] {
				if e.w > cur.last && !visited[e.to] {
					visited[e.to] = true
					ans[start]++
					q = append(q, State{v: e.to, last: e.w})
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
