package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAX_N = 2005

var (
	adj [MAX_N][]int
	par [MAX_N]int
	lvl [MAX_N]int
)

func dfs(u, p int, out *int) {
	lvl[u] = lvl[p] + 1
	par[u] = p
	if lvl[u] > lvl[*out] {
		*out = u
	}
	for _, nxt := range adj[u] {
		if nxt == p {
			continue
		}
		dfs(nxt, u, out)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	for i := 1; i <= n; i++ {
		adj[i] = adj[i][:0]
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	root := 1
	dfs(1, 0, &root)
	farth := root
	dfs(root, 0, &farth)
	var diam []int
	for u := farth; u != root; u = par[u] {
		diam = append(diam, u)
	}
	diam = append(diam, root)
	m := len(diam)
	var ops [][2]int
	if m%2 == 1 {
		c := diam[m/2]
		for i := 0; i <= m/2; i++ {
			ops = append(ops, [2]int{c, i})
		}
	} else {
		c1 := diam[m/2]
		c2 := diam[m/2-1]
		for i := m/2 - 1; i >= 0; i -= 2 {
			ops = append(ops, [2]int{c1, i})
			ops = append(ops, [2]int{c2, i})
		}
	}
	fmt.Fprintln(writer, len(ops))
	for _, pr := range ops {
		fmt.Fprintln(writer, pr[0], pr[1])
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var testc int
	fmt.Fscan(reader, &testc)
	for i := 0; i < testc; i++ {
		solve(reader, writer)
	}
}
