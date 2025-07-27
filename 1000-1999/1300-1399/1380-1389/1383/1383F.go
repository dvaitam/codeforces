package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1000000000

type Edge struct {
	to, rev, cap int
}

type Dinic struct {
	g     [][]Edge
	level []int
	iter  []int
}

func NewDinic(n int) *Dinic {
	g := make([][]Edge, n)
	level := make([]int, n)
	iter := make([]int, n)
	return &Dinic{g: g, level: level, iter: iter}
}

func (d *Dinic) AddEdge(from, to, cap int) {
	d.g[from] = append(d.g[from], Edge{to: to, rev: len(d.g[to]), cap: cap})
	d.g[to] = append(d.g[to], Edge{to: from, rev: len(d.g[from]) - 1, cap: 0})
}

func (d *Dinic) bfs(s int, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, len(d.g))
	d.level[s] = 0
	queue = append(queue, s)
	for q := 0; q < len(queue); q++ {
		v := queue[q]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(v, t, f int) int {
	if v == t {
		return f
	}
	for i := d.iter[v]; i < len(d.g[v]); i++ {
		d.iter[v] = i
		e := &d.g[v][i]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.g[e.to][e.rev].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k, q int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &q); err != nil {
		return
	}

	specials := make([][2]int, k)
	type EdgeInput struct{ u, v, w int }
	baseEdges := make([]EdgeInput, 0, m-k)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		if i < k {
			specials[i] = [2]int{u, v}
		} else {
			baseEdges = append(baseEdges, EdgeInput{u, v, w})
		}
	}

	size := 1 << k
	cost := make([]int, size)
	for mask := 0; mask < size; mask++ {
		din := NewDinic(n + 1)
		for _, e := range baseEdges {
			din.AddEdge(e.u, e.v, e.w)
		}
		for i := 0; i < k; i++ {
			u := specials[i][0]
			v := specials[i][1]
			if (mask>>i)&1 == 1 {
				din.AddEdge(1, u, INF)
				din.AddEdge(v, n, INF)
			} else {
				din.AddEdge(u, v, INF)
			}
		}
		cost[mask] = din.MaxFlow(1, n)
	}

	// mapping from single bit to index
	bitIndex := make([]int, size)
	for i := 0; i < k; i++ {
		bitIndex[1<<i] = i
	}
	weights := make([]int, k)
	sum := make([]int, size)

	for ; q > 0; q-- {
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &weights[i])
		}
		sum[0] = 0
		for mask := 1; mask < size; mask++ {
			lb := mask & -mask
			idx := bitIndex[lb]
			sum[mask] = sum[mask^lb] + weights[idx]
		}
		ans := INF
		for mask := 0; mask < size; mask++ {
			val := cost[mask] + sum[mask]
			if val < ans {
				ans = val
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
