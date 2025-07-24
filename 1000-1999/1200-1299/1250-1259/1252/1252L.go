package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	rev int
	cap int
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

func (d *Dinic) bfs(s int) {
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
	const inf = int(1 << 30)
	for {
		d.bfs(s)
		if d.level[t] < 0 {
			break
		}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, K int
	if _, err := fmt.Fscan(in, &N, &K); err != nil {
		return
	}

	A := make([]int, N+1)
	B := make([][]int, N+1)
	for i := 1; i <= N; i++ {
		var Ai, Mi int
		fmt.Fscan(in, &Ai, &Mi)
		A[i] = Ai
		B[i] = make([]int, Mi)
		for j := 0; j < Mi; j++ {
			fmt.Fscan(in, &B[i][j])
		}
	}

	C := make([]int, K)
	for i := 0; i < K; i++ {
		fmt.Fscan(in, &C[i])
	}

	// find cycle in the functional graph
	visited := make([]int, N+1)
	for i := range visited {
		visited[i] = -1
	}
	path := make([]int, 0, N)
	cur := 1
	for visited[cur] == -1 {
		visited[cur] = len(path)
		path = append(path, cur)
		cur = A[cur]
	}
	start := visited[cur]
	cycleNodes := path[start:]
	inCycle := make([]bool, N+1)
	for _, v := range cycleNodes {
		inCycle[v] = true
	}
	L := len(cycleNodes)

	// map materials to indices
	matIndex := map[int]int{}
	materials := make([]int, 0)
	addMat := func(x int) int {
		if idx, ok := matIndex[x]; ok {
			return idx
		}
		idx := len(materials)
		matIndex[x] = idx
		materials = append(materials, x)
		return idx
	}
	for i := 1; i <= N; i++ {
		for _, m := range B[i] {
			addMat(m)
		}
	}
	for _, c := range C {
		addMat(c)
	}
	M := len(materials)

	workers := make([][]int, M)
	for idx, c := range C {
		midx := matIndex[c]
		workers[midx] = append(workers[midx], idx)
	}

	src := 0
	gate := 1
	edgeStart := 2
	matStart := edgeStart + N
	sink := matStart + M
	dinic := NewDinic(sink + 1)

	if L > 0 {
		dinic.AddEdge(src, gate, L-1)
	}
	for i := 1; i <= N; i++ {
		node := edgeStart + i - 1
		if inCycle[i] {
			dinic.AddEdge(gate, node, 1)
		} else {
			dinic.AddEdge(src, node, 1)
		}
		for _, m := range B[i] {
			midx := matIndex[m]
			dinic.AddEdge(node, matStart+midx, 1)
		}
	}
	for m := 0; m < M; m++ {
		dinic.AddEdge(matStart+m, sink, len(workers[m]))
	}

	target := N - 1
	flow := dinic.MaxFlow(src, sink)
	if flow < target {
		fmt.Fprintln(out, -1)
		return
	}

	// determine assigned material for each edge
	assigned := make([]int, N+1)
	used := make([]bool, N+1)
	for i := 1; i <= N; i++ {
		node := edgeStart + i - 1
		for _, e := range dinic.g[node] {
			if e.to >= matStart && e.to < matStart+M && e.cap == 0 {
				assigned[i] = e.to - matStart
				used[i] = true
				break
			}
		}
	}

	workerAssign := make([][2]int, K)
	for i := 1; i <= N; i++ {
		if !used[i] {
			continue
		}
		midx := assigned[i]
		if len(workers[midx]) == 0 {
			continue
		}
		w := workers[midx][len(workers[midx])-1]
		workers[midx] = workers[midx][:len(workers[midx])-1]
		workerAssign[w] = [2]int{i, A[i]}
	}

	for i := 0; i < K; i++ {
		pair := workerAssign[i]
		fmt.Fprintf(out, "%d %d\n", pair[0], pair[1])
	}
}
