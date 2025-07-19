package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	p := make([]int, n+1)
	done := make([]bool, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	cnt := 0

	var dfs func(s, par int)
	dfs = func(s, par int) {
		c := -1
		for _, u := range adj[s] {
			if u == par {
				continue
			}
			c = u
			dfs(u, s)
		}
		if !done[s] && par != 0 {
			p[s], p[par] = p[par], p[s]
			done[s] = true
			done[par] = true
			cnt += 2
		}
		if !done[s] && c != -1 {
			p[s], p[c] = p[c], p[s]
			done[s] = true
			done[c] = true
			cnt += 2
		}
	}

	dfs(1, 0)
	fmt.Fprintln(writer, cnt)
	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, p[i])
	}
	fmt.Fprintln(writer)
}
