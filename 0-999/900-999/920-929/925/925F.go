package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to, rev int
	cap, flow float64
}

type EdgeInfo struct {
	u, v       int
	a, b, c, d float64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	edges := make([]EdgeInfo, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].a, &edges[i].b, &edges[i].c, &edges[i].d)
	}

	graph := make([][]Edge, n+2)
	for i := range graph {
		graph[i] = make([]Edge, 0, 10)
	}

	addEdge := func(u, v int, cap float64) {
		if u == v {
			return
		}
		graph[u] = append(graph[u], Edge{to: v, rev: len(graph[v]), cap: cap, flow: 0})
		graph[v] = append(graph[v], Edge{to: u, rev: len(graph[u]) - 1, cap: 0, flow: 0})
	}

	level := make([]int, n+2)
	ptr := make([]int, n+2)
	q := make([]int, 0, n+2)

	dinic := func(S, T int) float64 {
		var maxFlow float64
		for {
			for i := range level {
				level[i] = -1
			}
			level[S] = 0
			q = q[:0]
			q = append(q, S)
			head := 0
			for head < len(q) {
				u := q[head]
				head++
				for _, e := range graph[u] {
					if level[e.to] == -1 && e.cap-e.flow > 1e-9 {
						level[e.to] = level[u] + 1
						q = append(q, e.to)
					}
				}
			}
			if level[T] == -1 {
				break
			}

			for i := range ptr {
				ptr[i] = 0
			}

			var dfs func(u int, pushed float64) float64
			dfs = func(u int, pushed float64) float64 {
				if pushed < 1e-9 {
					return 0
				}
				if u == T {
					return pushed
				}
				for ; ptr[u] < len(graph[u]); ptr[u]++ {
					e := &graph[u][ptr[u]]
					tr := e.to
					if level[u]+1 != level[tr] || e.cap-e.flow < 1e-9 {
						continue
					}
					push := pushed
					if e.cap-e.flow < push {
						push = e.cap - e.flow
					}
					p := dfs(tr, push)
					if p < 1e-9 {
						continue
					}
					e.flow += p
					graph[tr][e.rev].flow -= p
					return p
				}
				return 0
			}

			for {
				pushed := dfs(S, 1e15)
				if pushed < 1e-9 {
					break
				}
				maxFlow += pushed
			}
		}
		return maxFlow
	}

	evaluate := func(t float64) float64 {
		for i := range graph {
			graph[i] = graph[i][:0]
		}

		D := make([]float64, n+2)
		for _, e := range edges {
			l := e.a*t + e.b
			r := e.c*t + e.d
			addEdge(e.u, e.v, r-l)
			D[e.v] += l
			D[e.u] -= l
		}

		var sumD float64
		for v := 1; v <= n; v++ {
			if D[v] > 0 {
				addEdge(0, v, D[v])
				sumD += D[v]
			} else if D[v] < 0 {
				addEdge(v, n+1, -D[v])
			}
		}

		return dinic(0, n+1) - sumD
	}

	L := 0.0
	R := 1.0
	for i := 0; i < 100; i++ {
		m1 := L + (R-L)/3.0
		m2 := R - (R-L)/3.0
		if evaluate(m1) < evaluate(m2) {
			L = m1
		} else {
			R = m2
		}
	}

	t_max := L
	max_g := evaluate(t_max)

	if max_g < -1e-7 {
		fmt.Printf("%.10f\n", 0.0)
		return
	}

	L_bin := 0.0
	R_bin := t_max
	for i := 0; i < 80; i++ {
		mid := L_bin + (R_bin-L_bin)/2.0
		if evaluate(mid) >= -1e-7 {
			R_bin = mid
		} else {
			L_bin = mid
		}
	}
	t_left := R_bin

	L_bin2 := t_max
	R_bin2 := 1.0
	for i := 0; i < 80; i++ {
		mid := L_bin2 + (R_bin2-L_bin2)/2.0
		if evaluate(mid) >= -1e-7 {
			L_bin2 = mid
		} else {
			R_bin2 = mid
		}
	}
	t_right := L_bin2

	fmt.Printf("%.10f\n", t_right-t_left)
}