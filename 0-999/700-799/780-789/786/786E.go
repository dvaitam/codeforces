package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

const (
	MAG = 100
)

type Edge struct{ to, cap, rev int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	// tree
	vv := make([][]int, n)
	Ex := make([]int, n)
	Ey := make([]int, n)
	ddd := make(map[[2]int]int)
	for i := 1; i < n; i++ {
		fmt.Fscan(in, &Ex[i], &Ey[i])
		Ex[i]--
		Ey[i]--
		vv[Ex[i]] = append(vv[Ex[i]], Ey[i])
		vv[Ey[i]] = append(vv[Ey[i]], Ex[i])
		ddd[[2]int{Ex[i], Ey[i]}] = i
		ddd[[2]int{Ey[i], Ex[i]}] = i
	}
	// LCA prep
	maxLG := 17
	par := make([][]int, maxLG)
	for i := range par {
		par[i] = make([]int, n)
	}
	depth := make([]int, n)
	cnt := make([]int, MAG)
	var dfs func(int, int)
	dfs = func(u, p int) {
		par[0][u] = p
		d := depth[u]
		cnt[d%MAG]++
		for _, w := range vv[u] {
			if w == p {
				continue
			}
			depth[w] = d + 1
			dfs(w, u)
		}
	}
	depth[0] = 0
	dfs(0, 0)
	for lg := 1; lg < maxLG; lg++ {
		for i := 0; i < n; i++ {
			par[lg][i] = par[lg-1][par[lg-1][i]]
		}
	}
	// choose G
	G := 0
	for i := 1; i < MAG; i++ {
		if cnt[i] < cnt[G] {
			G = i
		}
	}
	// lca func
	lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for lg := 0; diff > 0; lg++ {
			if diff&1 != 0 {
				u = par[lg][u]
			}
			diff >>= 1
		}
		if u == v {
			return u
		}
		for lg := maxLG - 1; lg >= 0; lg-- {
			if par[lg][u] != par[lg][v] {
				u = par[lg][u]
				v = par[lg][v]
			}
		}
		return par[0][u]
	}
	// read queries
	Qx := make([]int, m)
	Qy := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &Qx[i], &Qy[i])
		Qx[i]--
		Qy[i]--
	}
	// flow graph
	// base nodes: 0..m-1 queries, m..m+n-1 tree nodes
	S := m + n
	T := S + 1
	st := T
	// special nodes id, gs
	id := make([]int, n)
	gs := make([]int, n)
	for i := range id {
		id[i] = -1
	}
	// edges
	var Gph [][]Edge
	szGuess := m + n + n/MAG + 5
	Gph = make([][]Edge, szGuess)
	addEdge := func(u, v, c int) {
		Gph[u] = append(Gph[u], Edge{v, c, len(Gph[v])})
		Gph[v] = append(Gph[v], Edge{u, 0, len(Gph[u]) - 1})
	}
	// prepare ancestors for jumping MAG steps
	for i := 0; i < n; i++ {
		if depth[i]%MAG == G && depth[i] >= MAG {
			st++
			x := i
			for j := 0; j < MAG; j++ {
				addEdge(st, m+x, 1000000000)
				x = par[0][x]
			}
			id[i] = st
			gs[i] = x
		}
	}
	// connect queries
	for i := 0; i < m; i++ {
		z := lca(Qx[i], Qy[i])
		// from Qx[i] to z
		x := Qx[i]
		for x != z {
			if depth[x]%MAG == G && depth[x]-depth[z] >= MAG {
				addEdge(i, id[x], 1)
				x = gs[x]
				continue
			}
			addEdge(i, m+x, 1)
			x = par[0][x]
		}
		// from Qy[i]
		x = Qy[i]
		for x != z {
			if depth[x]%MAG == G && depth[x]-depth[z] >= MAG {
				addEdge(i, id[x], 1)
				x = gs[x]
				continue
			}
			addEdge(i, m+x, 1)
			x = par[0][x]
		}
	}
	// S to queries
	for i := 0; i < m; i++ {
		addEdge(S, i, 1)
	}
	// tree nodes to T
	for i := 0; i < n; i++ {
		addEdge(m+i, T, 1)
	}
	Nnodes := st + 1
	// Dinic
	level := make([]int, Nnodes)
	iter := make([]int, Nnodes)
	var bfs func()
	bfs = func() {
		for i := range level {
			level[i] = -1
		}
		queue := list.New()
		level[S] = 0
		queue.PushBack(S)
		for queue.Len() > 0 {
			u := queue.Remove(queue.Front()).(int)
			for _, e := range Gph[u] {
				if e.cap > 0 && level[e.to] < 0 {
					level[e.to] = level[u] + 1
					queue.PushBack(e.to)
				}
			}
		}
	}
	var dfsFlow func(int, int) int
	dfsFlow = func(u, f int) int {
		if u == T {
			return f
		}
		for i := iter[u]; i < len(Gph[u]); i++ {
			e := &Gph[u][i]
			if e.cap > 0 && level[u] < level[e.to] {
				d := dfsFlow(e.to, min(f, e.cap))
				if d > 0 {
					e.cap -= d
					Gph[e.to][e.rev].cap += d
					return d
				}
			}
			iter[u]++
		}
		return 0
	}
	flow := 0
	for {
		bfs()
		if level[T] < 0 {
			break
		}
		for i := range iter {
			iter[i] = 0
		}
		for {
			f := dfsFlow(S, 1<<60)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	// residual reachability
	vis := make([]bool, Nnodes)
	queue := list.New()
	queue.PushBack(S)
	vis[S] = true
	for queue.Len() > 0 {
		u := queue.Remove(queue.Front()).(int)
		for _, e := range Gph[u] {
			if e.cap > 0 && !vis[e.to] {
				vis[e.to] = true
				queue.PushBack(e.to)
			}
		}
	}
	// answer A and need
	var A []int
	need := make([]int, n)
	for i := 0; i < m; i++ {
		if !vis[i] {
			A = append(A, i+1)
		} else {
			need[Qx[i]]++
			need[Qy[i]]++
			z := lca(Qx[i], Qy[i])
			need[z] -= 2
		}
	}
	// accumulate need
	var dfs2 func(int, int)
	dfs2 = func(u, p int) {
		for _, w := range vv[u] {
			if w == p {
				continue
			}
			dfs2(w, u)
			need[u] += need[w]
		}
	}
	dfs2(0, -1)
	// B edges
	var B []int
	for i := 0; i < n; i++ {
		if need[i] > 0 {
			idx := ddd[[2]int{i, par[0][i]}]
			B = append(B, idx)
		}
	}
	// output
	fmt.Fprintln(out, flow)
	fmt.Fprint(out, len(A))
	for _, x := range A {
		fmt.Fprint(out, " ", x)
	}
	fmt.Fprintln(out)
	fmt.Fprint(out, len(B))
	for _, x := range B {
		fmt.Fprint(out, " ", x)
	}
	fmt.Fprintln(out)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
