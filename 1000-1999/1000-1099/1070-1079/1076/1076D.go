package main

import (
	"container/heap"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// State represents a node in the priority queue
type State struct {
	dist int64
	v    int
}

// PriorityQueue implements a min-heap for State
type PriorityQueue []State

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].dist == pq[j].dist {
		return pq[i].v < pq[j].v
	}
	return pq[i].dist < pq[j].dist
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func main() {
	// Read input
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	parts := strings.Fields(string(input))
	it := 0

	// Parse n, m, k
	n, _ := strconv.Atoi(parts[it])
	it++
	m, _ := strconv.Atoi(parts[it])
	it++
	k, _ := strconv.Atoi(parts[it])
	it++

	// Initialize graph
	g := make([][][3]int, n+1)
	for i := range g {
		g[i] = make([][3]int, 0)
	}
	for idx := 1; idx <= m; idx++ {
		x, _ := strconv.Atoi(parts[it])
		it++
		y, _ := strconv.Atoi(parts[it])
		it++
		w, _ := strconv.ParseInt(parts[it], 10, 64)
		it++
		g[x] = append(g[x], [3]int{y, int(w), idx})
		g[y] = append(g[y], [3]int{x, int(w), idx})
	}

	// Initialize distances and parent edges
	const INF = int64(1 << 60)
	dist := make([]int64, n+1)
	parEdge := make([]int, n+1)
	for i := range dist {
		dist[i] = INF
	}
	dist[1] = 0

	// Initialize priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, State{dist: 0, v: 1})

	// Dijkstra's algorithm
	for pq.Len() > 0 {
		state := heap.Pop(pq).(State)
		d, u := state.dist, state.v
		if d != dist[u] {
			continue
		}
		for _, edge := range g[u] {
			v, w, idx := edge[0], int64(edge[1]), edge[2]
			nd := d + w
			if nd < dist[v] {
				dist[v] = nd
				parEdge[v] = idx
				heap.Push(pq, State{dist: nd, v: v})
			}
		}
	}

	// Collect vertices with distances
	type Vert struct {
		dist int64
		v    int
	}
	verts := make([]Vert, 0, n-1)
	for v := 2; v <= n; v++ {
		verts = append(verts, Vert{dist[v], v})
	}
	sort.Slice(verts, func(i, j int) bool {
		if verts[i].dist == verts[j].dist {
			return verts[i].v < verts[j].v
		}
		return verts[i].dist < verts[j].dist
	})

	// Select up to k edges
	limit := min(k, n-1)
	ans := make([]int, 0, limit)
	for _, vert := range verts {
		if len(ans) == limit {
			break
		}
		ans = append(ans, parEdge[vert.v])
	}

	// Output results
	fmt.Println(len(ans))
	if len(ans) > 0 {
		for i, e := range ans {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(e)
		}
	}
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
