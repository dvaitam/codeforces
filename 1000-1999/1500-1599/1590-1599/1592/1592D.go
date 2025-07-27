package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

const LOG = 20

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type edge struct {
		to int
		w  int64
	}
	g := make([][]edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], edge{v, w})
		g[v] = append(g[v], edge{u, w})
	}

	parent := make([][LOG]int, n+1)
	upGcd := make([][LOG]int64, n+1)
	depth := make([]int, n+1)

	var dfs func(int, int)
	dfs = func(v, p int) {
		for _, e := range g[v] {
			if e.to == p {
				continue
			}
			depth[e.to] = depth[v] + 1
			parent[e.to][0] = v
			upGcd[e.to][0] = e.w
			for k := 1; k < LOG; k++ {
				parent[e.to][k] = parent[parent[e.to][k-1]][k-1]
				upGcd[e.to][k] = gcd(upGcd[e.to][k-1], upGcd[parent[e.to][k-1]][k-1])
			}
			dfs(e.to, v)
		}
	}
	dfs(1, 0)

	getGCD := func(u, v int) int64 {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		res := int64(0)
		diff := depth[u] - depth[v]
		for k := LOG - 1; k >= 0; k-- {
			if diff&(1<<k) != 0 {
				res = gcd(res, upGcd[u][k])
				u = parent[u][k]
			}
		}
		if u == v {
			return res
		}
		for k := LOG - 1; k >= 0; k-- {
			if parent[u][k] != parent[v][k] {
				res = gcd(res, upGcd[u][k])
				res = gcd(res, upGcd[v][k])
				u = parent[u][k]
				v = parent[v][k]
			}
		}
		res = gcd(res, upGcd[u][0])
		res = gcd(res, upGcd[v][0])
		return res
	}

	best := int64(0)
	ansU, ansV := 1, 1
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			val := getGCD(u, v)
			if val > best {
				best = val
				ansU, ansV = u, v
			}
		}
	}

	fmt.Fprintln(out, ansU, ansV)
}
