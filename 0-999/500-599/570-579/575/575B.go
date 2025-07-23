package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

type Edge struct{ to, id int }
type EdgeData struct {
	a, b int
	dir  int
}

var (
	n         int
	adj       [][]Edge
	edges     []EdgeData
	parent    [][]int
	depth     []int
	orient    []int
	upCount   []int64
	downCount []int64
)

func dfsInit(u, p int) {
	parent[0][u] = p
	for _, e := range adj[u] {
		v := e.to
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		ed := edges[e.id]
		if ed.dir == 0 {
			orient[v] = 0
		} else {
			if ed.a == u && ed.b == v {
				orient[v] = 1 // legal parent->child
			} else {
				orient[v] = 2 // legal child->parent
			}
		}
		dfsInit(v, u)
	}
}

func buildLCA() {
	for k := 1; k < len(parent); k++ {
		for i := 1; i <= n; i++ {
			parent[k][i] = parent[k-1][parent[k-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			a = parent[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := len(parent) - 1; k >= 0; k-- {
		if parent[k][a] != parent[k][b] {
			a = parent[k][a]
			b = parent[k][b]
		}
	}
	return parent[0][a]
}

func dfsAccum(u, p int) {
	for _, e := range adj[u] {
		v := e.to
		if v == p {
			continue
		}
		dfsAccum(v, u)
		upCount[u] += upCount[v]
		downCount[u] += downCount[v]
	}
}

func modPow(base, exp int64) int64 {
	res := int64(1)
	b := base % MOD
	for exp > 0 {
		if exp&1 == 1 {
			res = res * b % MOD
		}
		b = b * b % MOD
		exp >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	adj = make([][]Edge, n+1)
	edges = make([]EdgeData, n)
	for i := 1; i <= n-1; i++ {
		var a, b, x int
		fmt.Fscan(in, &a, &b, &x)
		edges[i] = EdgeData{a: a, b: b, dir: x}
		adj[a] = append(adj[a], Edge{to: b, id: i})
		adj[b] = append(adj[b], Edge{to: a, id: i})
	}
	maxLog := 0
	for (1 << maxLog) <= n {
		maxLog++
	}
	parent = make([][]int, maxLog)
	for i := range parent {
		parent[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	orient = make([]int, n+1)

	dfsInit(1, 1)
	buildLCA()

	var K int
	fmt.Fscan(in, &K)
	upCount = make([]int64, n+1)
	downCount = make([]int64, n+1)
	prev := 1
	for i := 0; i < K; i++ {
		var cur int
		fmt.Fscan(in, &cur)
		l := lca(prev, cur)
		upCount[prev]++
		upCount[l]--
		downCount[cur]++
		downCount[l]--
		prev = cur
	}
	dfsAccum(1, 1)
	ans := int64(0)
	for v := 2; v <= n; v++ {
		var m int64
		switch orient[v] {
		case 1:
			m = upCount[v]
		case 2:
			m = downCount[v]
		default:
			m = 0
		}
		if m > 0 {
			cost := modPow(2, m) - 1
			if cost < 0 {
				cost += MOD
			}
			ans += cost
			ans %= MOD
		}
	}
	fmt.Println(ans % MOD)
}
