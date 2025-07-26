package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

type Edge struct {
	to int
	w  int
}

var (
	n, q  int
	adj   [][]Edge
	up    [][]int
	mx    [][]int
	depth []int
	open  []bool
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dfs(v, p, w int) {
	up[v][0] = p
	mx[v][0] = w
	for i := 1; i < LOG; i++ {
		up[v][i] = up[up[v][i-1]][i-1]
		if up[v][i] == 0 {
			mx[v][i] = mx[v][i-1]
		} else {
			mx[v][i] = max(mx[v][i-1], mx[up[v][i-1]][i-1])
		}
	}
	for _, e := range adj[v] {
		if e.to == p {
			continue
		}
		depth[e.to] = depth[v] + 1
		dfs(e.to, v, e.w)
	}
}

func maxEdge(u, v int) int {
	if u == v {
		return 0
	}
	res := 0
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff>>i&1 == 1 {
			if mx[u][i] > res {
				res = mx[u][i]
			}
			u = up[u][i]
		}
	}
	if u == v {
		return res
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[u][i] != up[v][i] {
			if mx[u][i] > res {
				res = mx[u][i]
			}
			if mx[v][i] > res {
				res = mx[v][i]
			}
			u = up[u][i]
			v = up[v][i]
		}
	}
	if mx[u][0] > res {
		res = mx[u][0]
	}
	if mx[v][0] > res {
		res = mx[v][0]
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &q)
	adj = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
	}
	up = make([][]int, n+1)
	mx = make([][]int, n+1)
	for i := 0; i <= n; i++ {
		up[i] = make([]int, LOG)
		mx[i] = make([]int, LOG)
	}
	depth = make([]int, n+1)
	dfs(1, 1, 0)
	open = make([]bool, n+1)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			for i := l; i <= r; i++ {
				open[i] = true
			}
		} else if t == 2 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			for i := l; i <= r; i++ {
				open[i] = false
			}
		} else {
			var x int
			fmt.Fscan(reader, &x)
			ans := -1
			for i := 1; i <= n; i++ {
				if open[i] && i != x {
					w := maxEdge(x, i)
					if w > ans {
						ans = w
					}
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
