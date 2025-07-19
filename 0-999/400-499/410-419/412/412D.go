package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	n, m int
	g    [][]int
	vis  []bool
	out  *bufio.Writer
)

func dfs(v int) {
	vis[v] = true
	for _, u := range g[v] {
		if !vis[u] {
			dfs(u)
		}
	}
	out.WriteString(strconv.Itoa(v + 1))
	out.WriteByte(' ')
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g = make([][]int, n)
	vis = make([]bool, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
	}
	for v := 0; v < n; v++ {
		if !vis[v] {
			dfs(v)
		}
	}
	out.WriteByte('\n')
}
