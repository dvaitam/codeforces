package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to  int
	rev int
	cap int64
}

type dinic struct {
	g     [][]edge
	level []int
	it    []int
	s, t  int
}

func newDinic(n, s, t int) *dinic {
	g := make([][]edge, n)
	return &dinic{
		g:     g,
		level: make([]int, n),
		it:    make([]int, n),
		s:     s,
		t:     t,
	}
}

func (d *dinic) addEdge(u, v int, c int64) {
	d.g[u] = append(d.g[u], edge{to: v, rev: len(d.g[v]), cap: c})
	d.g[v] = append(d.g[v], edge{to: u, rev: len(d.g[u]) - 1, cap: 0})
}

func (d *dinic) addUndirected(u, v int, c int64) {
	d.addEdge(u, v, c)
	d.addEdge(v, u, c)
}

func (d *dinic) bfs() bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, len(d.level))
	d.level[d.s] = 0
	queue = append(queue, d.s)
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, e := range d.g[u] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[u] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.level[d.t] >= 0
}

func (d *dinic) dfs(u int, f int64) int64 {
	if u == d.t {
		return f
	}
	for ; d.it[u] < len(d.g[u]); d.it[u]++ {
		e := &d.g[u][d.it[u]]
		if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
			ret := d.dfs(e.to, min64(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.g[e.to][e.rev].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *dinic) maxFlow() int64 {
	var flow int64
	for d.bfs() {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(d.s, 1<<60)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func node(v, color int) int {
	return (v-1)*2 + color
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var a, b int
	fmt.Fscan(in, &a, &b)

	totalNodes := 2*n + 2
	s := 2 * n
	t := 2*n + 1
	const inf int64 = 1 << 62

	d := newDinic(totalNodes, s, t)

	// Fix castles: a is color 0 (A), b is color 1 (B).
	d.addEdge(s, node(a, 0), inf)
	d.addEdge(node(a, 1), t, inf)
	d.addEdge(s, node(b, 1), inf)
	d.addEdge(node(b, 0), t, inf)

	for i := 0; i < m; i++ {
		var u, v int
		var l int64
		fmt.Fscan(in, &u, &v, &l)
		u0, u1 := node(u, 0), node(u, 1)
		v0, v1 := node(v, 0), node(v, 1)
		// Two edges in the doubled graph reflect the cost table.
		d.addUndirected(u0, v1, l)
		d.addUndirected(u1, v0, l)
	}

	ans := d.maxFlow()

	// Recover partition by reachability in residual graph.
	reach := make([]bool, totalNodes)
	queue := []int{s}
	reach[s] = true
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, e := range d.g[u] {
			if e.cap > 0 && !reach[e.to] {
				reach[e.to] = true
				queue = append(queue, e.to)
			}
		}
	}

	result := make([]byte, n)
	for v := 1; v <= n; v++ {
		id0 := node(v, 0)
		id1 := node(v, 1)
		in0 := reach[id0]
		in1 := reach[id1]
		switch {
		case in0 && !in1:
			result[v-1] = 'A'
		case in1 && !in0:
			result[v-1] = 'B'
		default:
			result[v-1] = 'C'
		}
	}

	fmt.Fprintln(out, ans)
	fmt.Fprintln(out, string(result))
}
