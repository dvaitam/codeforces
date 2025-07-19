package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = 1000000000

// Edge represents a flow edge with two-dimensional cost
type Edge struct {
	to, cap, rev, c1, c2 int
}

var (
	g        [][]Edge
	dist1    []int
	dist2    []int
	inq      []bool
	prevV    []int
	prevE    []int
	n, x, y  int
	a        []int
	g1, g2   [][]int
	r1, r2   []int
	src, tgt int
	wts      int
)

func addEdge(u, v, cap, c1, c2 int) {
	g[u] = append(g[u], Edge{v, cap, len(g[v]), c1, c2})
	g[v] = append(g[v], Edge{u, 0, len(g[u]) - 1, -c1, -c2})
}

// spfa finds shortest path in lex order of (dist1,dist2)
func spfa() bool {
	tot := len(g)
	for i := 0; i < tot; i++ {
		dist1[i] = INF
		dist2[i] = INF
		inq[i] = false
	}
	dist1[src], dist2[src] = 0, 0
	q := make([]int, 0, tot)
	q = append(q, src)
	inq[src] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		inq[u] = false
		for ei, e := range g[u] {
			if e.cap > 0 {
				d1 := dist1[u] + e.c1
				d2 := dist2[u] + e.c2
				if dist1[e.to] > d1 || (dist1[e.to] == d1 && dist2[e.to] > d2) {
					dist1[e.to], dist2[e.to] = d1, d2
					prevV[e.to], prevE[e.to] = u, ei
					if !inq[e.to] {
						inq[e.to] = true
						q = append(q, e.to)
					}
				}
			}
		}
	}
	// return true if cost < (0,0)
	return dist1[tgt] < 0 || (dist1[tgt] == 0 && dist2[tgt] < 0)
}

func minCostMaxFlow() (int, int) {
	flow, cost1, cost2 := 0, 0, 0
	for spfa() {
		// find minimum capacity
		f := INF
		for v := tgt; v != src; v = prevV[v] {
			e := g[prevV[v]][prevE[v]]
			if e.cap < f {
				f = e.cap
			}
		}
		// apply flow
		for v := tgt; v != src; v = prevV[v] {
			e := &g[prevV[v]][prevE[v]]
			e.cap -= f
			cost1 += f * e.c1
			cost2 += f * e.c2
			g[v][e.rev].cap += f
		}
		flow += f
	}
	return cost1, cost2
}

func dfs1(u, p, pre int) {
	if pre >= 0 {
		if r1[u] > 0 {
			addEdge(pre, n+u, r1[u], -1, 0)
			addEdge(n+u, u, 1, 0, -a[u])
		} else if pre != src {
			addEdge(pre, u, 1, 0, -a[u])
		}
	}
	for _, v := range g1[u] {
		if v == p {
			continue
		}
		nxt := pre
		if r1[u] > 0 {
			nxt = n + u
		}
		dfs1(v, u, nxt)
	}
}

func dfs2(u, p, pre int) {
	if pre >= 0 {
		if r2[u] > 0 {
			addEdge(2*n+u, pre, r2[u], -1, 0)
			addEdge(u, 2*n+u, 1, 0, 0)
		} else if pre != tgt {
			addEdge(u, pre, 1, 0, 0)
		}
	}
	for _, v := range g2[u] {
		if v == p {
			continue
		}
		nxt := pre
		if r2[u] > 0 {
			nxt = 2*n + u
		}
		dfs2(v, u, nxt)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &x, &y)
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	g1 = make([][]int, n)
	g2 = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g1[u] = append(g1[u], v)
		g1[v] = append(g1[v], u)
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g2[u] = append(g2[u], v)
		g2[v] = append(g2[v], u)
	}
	r1 = make([]int, n)
	r2 = make([]int, n)
	var q1, q2 int
	fmt.Fscan(in, &q1)
	for i := 0; i < q1; i++ {
		var k, num int
		fmt.Fscan(in, &k, &num)
		k--
		if r1[k] > 0 {
			fmt.Println(-1)
			return
		}
		r1[k] = num
		wts += num
	}
	fmt.Fscan(in, &q2)
	for i := 0; i < q2; i++ {
		var k, num int
		fmt.Fscan(in, &k, &num)
		k--
		if r2[k] > 0 {
			fmt.Println(-1)
			return
		}
		r2[k] = num
		wts += num
	}
	// build flow graph
	src = 3 * n
	tgt = 3*n + 1
	tot := 3*n + 2
	g = make([][]Edge, tot)
	dist1 = make([]int, tot)
	dist2 = make([]int, tot)
	inq = make([]bool, tot)
	prevV = make([]int, tot)
	prevE = make([]int, tot)
	// dfs trees
	dfs1(x-1, -1, src)
	dfs2(y-1, -1, tgt)
	// compute min cost max flow
	cost1, cost2 := minCostMaxFlow()
	if cost1 != -wts {
		fmt.Println(-1)
	} else {
		fmt.Println(-cost2)
	}
}
