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
	n     int
	g     [][]Edge
	level []int
	iter  []int
}

func NewDinic(n int) *Dinic {
	g := make([][]Edge, n)
	level := make([]int, n)
	iter := make([]int, n)
	return &Dinic{n: n, g: g, level: level, iter: iter}
}

func (d *Dinic) AddEdge(from, to int, cap int64) {
	d.g[from] = append(d.g[from], Edge{to: to, rev: len(d.g[to]), cap: cap})
	d.g[to] = append(d.g[to], Edge{to: from, rev: len(d.g[from]) - 1, cap: 0})
}

func (d *Dinic) bfs(s int) {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, d.n)
	d.level[s] = 0
	queue = append(queue, s)
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
}

func (d *Dinic) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for ; d.iter[v] < len(d.g[v]); d.iter[v]++ {
		e := &d.g[v][d.iter[v]]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min64(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.g[e.to][e.rev].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
	var flow int64
	const INF int64 = 1 << 60
	for {
		d.bfs(s)
		if d.level[t] < 0 {
			break
		}
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, INF)
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

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	x := make([]int64, n-1)
	y := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}

	type CE struct {
		a int
		b int
		c int64
	}
	cross := make([]CE, m)
	for i := 0; i < m; i++ {
		var a, b int
		var c int64
		fmt.Fscan(reader, &a, &b, &c)
		cross[i] = CE{a - 1, b - 1, c}
	}

	build := func() int64 {
		N := 2 * n
		d := NewDinic(N)
		for i := 0; i < n-1; i++ {
			d.AddEdge(i, i+1, x[i])
		}
		for i := 0; i < n-1; i++ {
			d.AddEdge(n+i, n+i+1, y[i])
		}
		for _, e := range cross {
			d.AddEdge(e.a, n+e.b, e.c)
		}
		return d.MaxFlow(0, N-1)
	}

	res := build()
	fmt.Fprint(writer, res)

	for i := 0; i < q; i++ {
		var v int
		var w int64
		fmt.Fscan(reader, &v, &w)
		x[v-1] = w
		res = build()
		fmt.Fprint(writer, " ", res)
	}
}
