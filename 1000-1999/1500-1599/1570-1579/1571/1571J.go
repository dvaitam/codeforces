package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	rev int
	cap int64
}

type Dinic struct {
	g     [][]Edge
	level []int
	it    []int
}

func NewDinic(n int) *Dinic {
	return &Dinic{g: make([][]Edge, n), level: make([]int, n), it: make([]int, n)}
}

func (d *Dinic) AddEdge(from, to int, cap int64) {
	d.g[from] = append(d.g[from], Edge{to, len(d.g[to]), cap})
	d.g[to] = append(d.g[to], Edge{from, len(d.g[from]) - 1, 0})
}

func (d *Dinic) AddBiEdge(u, v int, cap int64) {
	d.AddEdge(u, v, cap)
	d.AddEdge(v, u, cap)
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0)
	d.level[s] = 0
	q = append(q, s)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for ; d.it[v] < len(d.g[v]); d.it[v]++ {
		e := &d.g[v][d.it[v]]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			dflow := d.dfs(e.to, t, min64(f, e.cap))
			if dflow > 0 {
				e.cap -= dflow
				rev := &d.g[e.to][e.rev]
				rev.cap += dflow
				return dflow
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
	flow := int64(0)
	for d.bfs(s, t) {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<60)
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	c := make([]int64, n+1)
	for i := 1; i <= n-1; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n-1; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &c[i])
	}
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		// build graph for cities [l, r]
		m := 2 * (r - l + 1)
		din := NewDinic(m + 2)
		s := 0
		t := m + 1
		start := func(i int, up bool) int { // return index of node
			if up {
				return 1 + 2*(i-l)
			}
			return 1 + 2*(i-l) + 1
		}
		din.AddEdge(s, start(l, true), 1<<60)
		for i := l; i <= r; i++ {
			din.AddBiEdge(start(i, true), start(i, false), c[i])
		}
		for i := l; i < r; i++ {
			din.AddEdge(start(i, true), start(i+1, true), a[i])
			din.AddEdge(start(i, false), start(i+1, false), b[i])
		}
		din.AddEdge(start(r, true), t, 1<<60)
		flow := din.MaxFlow(s, t)
		fmt.Fprintln(writer, flow)
	}
}
