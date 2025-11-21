package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type edge struct {
	to  int
	rev int
	cap int64
}

type dinic struct {
	graph [][]edge
	level []int
	iter  []int
}

func newDinic(n int) *dinic {
	return &dinic{
		graph: make([][]edge, n),
		level: make([]int, n),
		iter:  make([]int, n),
	}
}

func (d *dinic) addEdge(from, to int, cap int64) {
	f := edge{to: to, rev: len(d.graph[to]), cap: cap}
	r := edge{to: from, rev: len(d.graph[from]), cap: 0}
	d.graph[from] = append(d.graph[from], f)
	d.graph[to] = append(d.graph[to], r)
}

func (d *dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, len(d.graph))
	head := 0
	d.level[s] = 0
	queue = append(queue, s)
	for head < len(queue) {
		v := queue[head]
		head++
		for _, e := range d.graph[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *dinic) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for i := d.iter[v]; i < len(d.graph[v]); i++ {
		d.iter[v] = i
		e := &d.graph[v][i]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			minCap := f
			if e.cap < minCap {
				minCap = e.cap
			}
			dFlow := d.dfs(e.to, t, minCap)
			if dFlow > 0 {
				e.cap -= dFlow
				rev := &d.graph[e.to][e.rev]
				rev.cap += dFlow
				return dFlow
			}
		}
	}
	d.iter[v] = len(d.graph[v])
	return 0
}

func (d *dinic) maxFlow(s, t int) int64 {
	var flow int64
	for d.bfs(s, t) {
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, inf)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var h, b int64
		fmt.Fscan(in, &n, &h, &b)
		grid := make([]string, 2)
		for i := 0; i < 2; i++ {
			fmt.Fscan(in, &grid[i])
		}

		cells := 2 * n
		totalNodes := cells*2 + 1
		source := totalNodes - 1
		d := newDinic(totalNodes)

		var sink int
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				idx := r*n + c
				inNode := idx * 2
				outNode := inNode + 1
				ch := grid[r][c]
				var cost int64
				switch ch {
				case 'S':
					cost = inf
					sink = outNode
				case 'W':
					cost = h
				default:
					cost = b
				}
				d.addEdge(inNode, outNode, cost)
				if ch == 'W' {
					d.addEdge(source, inNode, inf)
				}
			}
		}

		dir := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				idx := r*n + c
				outNode := idx*2 + 1
				for _, dxy := range dir {
					nr := r + dxy[0]
					nc := c + dxy[1]
					if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
						continue
					}
					nIdx := nr*n + nc
					nIn := nIdx * 2
					d.addEdge(outNode, nIn, inf)
				}
			}
		}

		result := d.maxFlow(source, sink)
		fmt.Fprintln(out, result)
	}
}
