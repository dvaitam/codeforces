package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type edge struct{ to int }
	g := make([][]edge, n+1)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = [2]int{u, v}
		g[u] = append(g[u], edge{v})
		g[v] = append(g[v], edge{u})
	}

	parent := make([]int, n+1)
	size := make([]int, n+1)

	var dfs func(int, int)
	dfs = func(v, p int) {
		parent[v] = p
		size[v] = 1
		for _, e := range g[v] {
			if e.to == p {
				continue
			}
			dfs(e.to, v)
			size[v] += size[e.to]
		}
	}
	dfs(1, 0)

	var ans int64 = 1
	for _, e := range edges {
		u, v := e[0], e[1]
		if parent[u] == v {
			u, v = v, u
		}
		s := size[v]
		if parent[v] != u {
			s = size[u]
		}
		other := n - s
		ways := int64(s) * int64(other) % mod
		ans = ans * ways % mod
	}

	fmt.Fprintln(out, ans)
}
