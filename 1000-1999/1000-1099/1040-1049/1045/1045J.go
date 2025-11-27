package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	c  byte
}

const LOG = 20

var (
	n        int
	g        [][]edge
	parent   [LOG][]int
	depth    []int
	charFrom []byte
)

func dfs(u, p int, c byte) {
	parent[0][u] = p
	charFrom[u] = c
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		depth[e.to] = depth[u] + 1
		dfs(e.to, u, e.c)
	}
}

func buildLCA() {
	for k := 1; k < LOG; k++ {
		for i := 0; i < n; i++ {
			par := parent[k-1][i]
			if par != -1 {
				parent[k][i] = parent[k-1][par]
			} else {
				parent[k][i] = -1
			}
		}
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	// lift u
	diff := depth[u] - depth[v]
	for k := LOG - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			u = parent[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := LOG - 1; k >= 0; k-- {
		if parent[k][u] != -1 && parent[k][u] != parent[k][v] {
			u = parent[k][u]
			v = parent[k][v]
		}
	}
	return parent[0][u]
}

func pathString(u, v int) string {
	l := lca(u, v)
	var up []byte
	for u != l {
		up = append(up, charFrom[u])
		u = parent[0][u]
	}
	var down []byte
	for v != l {
		down = append(down, charFrom[v])
		v = parent[0][v]
	}
	// reverse down
	for i, j := 0, len(down)-1; i < j; i, j = i+1, j-1 {
		down[i], down[j] = down[j], down[i]
	}
	return string(append(up, down...))
}

func countOccurrences(text, pat string) int {
	if len(pat) == 0 || len(text) < len(pat) {
		return 0
	}
	cnt := 0
	m := len(pat)
	for i := 0; i+m <= len(text); i++ {
		if text[i:i+m] == pat {
			cnt++
		}
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]edge, n)
	charFrom = make([]byte, n)
	depth = make([]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		var c string
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		ch := c[0]
		g[a] = append(g[a], edge{b, ch})
		g[b] = append(g[b], edge{a, ch})
	}
	for k := 0; k < LOG; k++ {
		parent[k] = make([]int, n)
		for i := range parent[k] {
			parent[k][i] = -1
		}
	}
	dfs(0, -1, 0)
	buildLCA()

	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < q; i++ {
		var u, v int
		var s string
		fmt.Fscan(in, &u, &v, &s)
		u--
		v--
		path := pathString(u, v)
		fmt.Fprintln(out, countOccurrences(path, s))
	}
}
