package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	// build parent and order starting from 1
	parent := make([]int, n+1)
	weightPar := make([]int64, n+1)
	order := make([]int, 0, n)
	queue := make([]int, 0, n)
	queue = append(queue, 1)
	parent[1] = 0
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		order = append(order, v)
		for _, e := range g[v] {
			if e.to == parent[v] {
				continue
			}
			parent[e.to] = v
			weightPar[e.to] = e.w
			queue = append(queue, e.to)
		}
	}

	// compute down values
	down := make([]int64, n+1)
	best1 := make([]int64, n+1)
	best2 := make([]int64, n+1)
	bestChild := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		for _, e := range g[v] {
			if e.to == parent[v] {
				continue
			}
			cand := down[e.to] + e.w
			if cand > best1[v] {
				best2[v] = best1[v]
				best1[v] = cand
				bestChild[v] = e.to
			} else if cand > best2[v] {
				best2[v] = cand
			}
		}
		down[v] = best1[v]
	}

	// compute up values
	up := make([]int64, n+1)
	for _, v := range order {
		for _, e := range g[v] {
			if e.to == parent[v] {
				continue
			}
			use := best1[v]
			if bestChild[v] == e.to {
				use = best2[v]
			}
			if up[v] > use {
				up[e.to] = e.w + up[v]
			} else {
				up[e.to] = e.w + use
			}
		}
	}

	// helper BFS to compute distances from a start node
	bfs := func(start int) ([]int64, []int, []int64) {
		dist := make([]int64, n+1)
		for i := range dist {
			dist[i] = -1
		}
		par := make([]int, n+1)
		wpar := make([]int64, n+1)
		q := make([]int, 0, n)
		q = append(q, start)
		dist[start] = 0
		par[start] = 0
		for head := 0; head < len(q); head++ {
			v := q[head]
			for _, e := range g[v] {
				if dist[e.to] != -1 {
					continue
				}
				dist[e.to] = dist[v] + e.w
				par[e.to] = v
				wpar[e.to] = e.w
				q = append(q, e.to)
			}
		}
		return dist, par, wpar
	}

	// get diameter endpoints A and B
	dist1, _, _ := bfs(1)
	A := 1
	for i := 1; i <= n; i++ {
		if dist1[i] > dist1[A] {
			A = i
		}
	}
	distA, _, wParA := bfs(A)
	B := A
	for i := 1; i <= n; i++ {
		if distA[i] > distA[B] {
			B = i
		}
	}
	distB, _, wParB := bfs(B)
	D0 := distA[B]

	delta := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if distA[i]+distB[i] == D0 {
			if i != A {
				delta[i] += wParA[i]
			}
			if i != B {
				delta[i] += wParB[i]
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var x int
		var k int64
		fmt.Fscan(in, &x, &k)
		bestOne := int64(0)
		bestTwo := int64(0)
		for _, e := range g[x] {
			var rest int64
			if parent[e.to] == x { // child
				rest = down[e.to]
			} else { // parent
				rest = up[x] - e.w
			}
			val := rest + e.w*k
			if val > bestOne {
				bestTwo = bestOne
				bestOne = val
			} else if val > bestTwo {
				bestTwo = val
			}
		}
		ans := bestOne
		if len(g[x]) >= 2 {
			if bestOne+bestTwo > ans {
				ans = bestOne + bestTwo
			}
		}
		dab := D0
		if delta[x] > 0 {
			dab = D0 + (k-1)*delta[x]
		}
		if dab > ans {
			ans = dab
		}
		fmt.Fprintln(out, ans)
	}
}
