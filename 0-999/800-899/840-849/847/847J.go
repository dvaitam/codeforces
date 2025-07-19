package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Edge struct {
	to, cap, rev int
}

var (
	n, m         int
	g            [][]Edge
	level, iter  []int
	edn          []int
	edm          []int
	ed           [][2]int
	source, sink int
)

func addEdge(u, v, c int) {
	g[u] = append(g[u], Edge{v, c, len(g[v])})
	g[v] = append(g[v], Edge{u, 0, len(g[u]) - 1})
}

func bfs() bool {
	for i := range level {
		level[i] = -1
		iter[i] = 0
	}
	queue := make([]int, 0, len(g))
	level[source] = 0
	queue = append(queue, source)
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, e := range g[v] {
			if e.cap > 0 && level[e.to] < 0 {
				level[e.to] = level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return level[sink] >= 0
}

func dfs(v, f int) int {
	if v == sink {
		return f
	}
	for i := iter[v]; i < len(g[v]); i++ {
		e := &g[v][i]
		if e.cap > 0 && level[v] < level[e.to] {
			d := dfs(e.to, min(f, e.cap))
			if d > 0 {
				e.cap -= d
				g[e.to][e.rev].cap += d
				return d
			}
		}
		iter[v]++
	}
	return 0
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
	readInt := func() int {
		s, _ := reader.ReadBytes(' ')
		v, _ := strconv.Atoi(string(s[:len(s)-1]))
		return v
	}
	// read n, m
	fmt.Fscan(reader, &n, &m)
	ed = make([][2]int, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &ed[i][0], &ed[i][1])
	}
	V := m + n + 2
	source = 0
	sink = m + n + 1
	g = make([][]Edge, V)
	level = make([]int, V)
	iter = make([]int, V)
	edn = make([]int, m+1)
	edm = make([]int, n+1)
	// build initial graph
	for i := 1; i <= m; i++ {
		addEdge(source, i, 1)
		addEdge(i, m+ed[i][0], 1)
		edn[i] = len(g[i]) - 1
		addEdge(i, m+ed[i][1], 1)
	}
	for i := 1; i <= n; i++ {
		addEdge(m+i, sink, 0)
		edm[i] = len(g[m+i]) - 1
	}
	flow := 0
	k := 0
	for flow < m {
		if !bfs() {
			k++
			for i := 1; i <= n; i++ {
				g[m+i][edm[i]].cap++
			}
			continue
		}
		for {
			f := dfs(source, 1<<60)
			if f <= 0 {
				break
			}
			flow += f
		}
	}
	// output
	fmt.Fprintln(writer, k)
	for i := 1; i <= m; i++ {
		if g[i][edn[i]].cap == 0 {
			fmt.Fprintf(writer, "%d %d\n", ed[i][0], ed[i][1])
		} else {
			fmt.Fprintf(writer, "%d %d\n", ed[i][1], ed[i][0])
		}
	}
}
