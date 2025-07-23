package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to  int
	rev int
	cap int
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

func (d *Dinic) AddEdge(from, to, cap int) {
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
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
}

func (d *Dinic) dfs(v, t, f int) int {
	if v == t {
		return f
	}
	for i := d.iter[v]; i < len(d.g[v]); i++ {
		e := &d.g[v][i]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.g[e.to][e.rev].cap += ret
				return ret
			}
		}
		d.iter[v]++
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for {
		d.bfs(s)
		if d.level[t] < 0 {
			break
		}
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<30)
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

func countRem(x, r int) int {
	full := x / 5
	rem := x % 5
	if r == 0 {
		return full
	}
	if r <= rem {
		return full + 1
	}
	return full
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, b, q int
	if _, err := fmt.Fscan(reader, &n, &b, &q); err != nil {
		return
	}

	type pair struct{ x, y int }
	hints := make([]pair, 0, q+2)
	for i := 0; i < q; i++ {
		var up, cnt int
		fmt.Fscan(reader, &up, &cnt)
		if up > b {
			up = b
		}
		if cnt > n {
			cnt = n
		}
		hints = append(hints, pair{up, cnt})
	}
	hints = append(hints, pair{b, n})

	sort.Slice(hints, func(i, j int) bool { return hints[i].x < hints[j].x })

	points := []pair{{0, 0}}
	for _, h := range hints {
		if h.x == points[len(points)-1].x {
			if h.y != points[len(points)-1].y {
				fmt.Fprint(writer, "unfair")
				return
			}
			continue
		}
		if h.y < points[len(points)-1].y || h.y > h.x {
			fmt.Fprint(writer, "unfair")
			return
		}
		if h.y-points[len(points)-1].y > h.x-points[len(points)-1].x {
			fmt.Fprint(writer, "unfair")
			return
		}
		points = append(points, h)
	}

	k := len(points) - 1
	segNeed := make([]int, k)
	segCount := make([][5]int, k)
	for i := 1; i <= k; i++ {
		L := points[i-1].x
		R := points[i].x
		need := points[i].y - points[i-1].y
		segNeed[i-1] = need
		for r := 0; r < 5; r++ {
			segCount[i-1][r] = countRem(R, r) - countRem(L, r)
		}
	}

	totalNodes := 1 + 5 + k + 1
	S := 0
	T := totalNodes - 1
	d := NewDinic(totalNodes)

	perClass := n / 5
	for r := 0; r < 5; r++ {
		d.AddEdge(S, 1+r, perClass)
	}
	for i := 0; i < k; i++ {
		idx := 1 + 5 + i
		d.AddEdge(idx, T, segNeed[i])
		for r := 0; r < 5; r++ {
			cap := segCount[i][r]
			if cap > 0 {
				d.AddEdge(1+r, idx, cap)
			}
		}
	}

	flow := d.MaxFlow(S, T)
	if flow == n {
		fmt.Fprint(writer, "fair")
	} else {
		fmt.Fprint(writer, "unfair")
	}
}
