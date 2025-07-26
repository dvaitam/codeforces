package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct{ to, id int }

var (
	n, m      int
	g         [][]edge
	up        [][]int
	depth     []int
	LOG       int
	tag       []int
	cnt       []int
	need      []bool
	mandatory []bool
)

func addEdge(u, v, id int) {
	g[u] = append(g[u], edge{v, id})
	g[v] = append(g[v], edge{u, id})
}

func dfsInit(u, p int) {
	up[0][u] = p
	for i := 1; i < LOG; i++ {
		up[i][u] = up[i-1][up[i-1][u]]
	}
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		depth[e.to] = depth[u] + 1
		dfsInit(e.to, u)
	}
}

func jump(u, k int) int {
	for i := 0; i < LOG; i++ {
		if k&1 == 1 {
			u = up[i][u]
		}
		k >>= 1
		if u == 0 {
			return 0
		}
	}
	return u
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := 0; i < LOG; i++ {
		if diff>>i&1 == 1 {
			a = up[i][a]
		}
	}
	if a == b {
		return a
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][a] != up[i][b] {
			a = up[i][a]
			b = up[i][b]
		}
	}
	return up[0][a]
}

func nextOnPath(a, b int) int {
	l := lca(a, b)
	if l != a {
		return up[0][a]
	}
	return jump(b, depth[b]-depth[a]-1)
}

func dfsCount(u, p int) {
	cnt[u] += tag[u]
	for _, e := range g[u] {
		v := e.to
		if v == p {
			continue
		}
		dfsCount(v, u)
		cnt[u] += cnt[v]
		if cnt[v] > 0 {
			need[e.id] = true
		}
	}
}

const INF int = 1 << 30

var dp0, dp1 []int

func dfsDP(u, p int) {
	dp1[u] = 1
	if mandatory[u] {
		dp0[u] = INF
	} else {
		dp0[u] = 0
	}
	for _, e := range g[u] {
		v := e.to
		if v == p {
			continue
		}
		dfsDP(v, u)
		if need[e.id] {
			dp1[u] += min(dp0[v], dp1[v])
			dp0[u] += dp1[v]
		} else {
			dp1[u] += min(dp0[v], dp1[v])
			dp0[u] += min(dp0[v], dp1[v])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m)
	g = make([][]edge, n+1)
	up = make([][]int, 20)
	for i := 0; i < 20; i++ {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		addEdge(p, i, i-1)
	}
	LOG = 20
	dfsInit(1, 0)
	tag = make([]int, n+1)
	cnt = make([]int, n+1)
	need = make([]bool, n)
	mandatory = make([]bool, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		l := lca(x, y)
		dist := depth[x] + depth[y] - 2*depth[l]
		if dist == 1 {
			fmt.Println(-1)
			return
		}
		nx := nextOnPath(x, y)
		ny := nextOnPath(y, x)
		if nx == ny {
			mandatory[nx] = true
		} else {
			l2 := lca(nx, ny)
			tag[nx]++
			tag[ny]++
			tag[l2] -= 2
		}
	}
	dfsCount(1, 0)
	dp0 = make([]int, n+1)
	dp1 = make([]int, n+1)
	dfsDP(1, 0)
	ans := min(dp0[1], dp1[1])
	fmt.Println(ans)
}
