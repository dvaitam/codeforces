package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Edge represents moving a free cell from the source to the destination at some cost.
type Edge struct {
	to   int
	cost int64
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var p, q int64
	fmt.Fscan(in, &p, &q)

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	idx := func(r, c int) int { return r*m + c }
	valid := func(r, c int) bool {
		return r >= 0 && r < n && c >= 0 && c < m && grid[r][c] != '#'
	}

	adj := make([][]Edge, n*m)
	dotCount := 0

	addEdge := func(fr, fc, toR, toC int, cost int64) {
		if !valid(fr, fc) || !valid(toR, toC) {
			return
		}
		from := idx(fr, fc)
		to := idx(toR, toC)
		adj[from] = append(adj[from], Edge{to: to, cost: cost})
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			switch grid[i][j] {
			case '.':
				dotCount++
			case 'L':
				addEdge(i, j-1, i, j+1, q) // shift left
				addEdge(i, j+2, i, j, q)   // shift right
				addEdge(i-1, j, i, j+1, p) // rotate around L
				addEdge(i+1, j, i, j+1, p)
				addEdge(i-1, j+1, i, j, p) // rotate around R
				addEdge(i+1, j+1, i, j, p)
			case 'U':
				addEdge(i-1, j, i+1, j, q) // shift up
				addEdge(i+2, j, i, j, q)   // shift down
				addEdge(i, j-1, i+1, j, p) // rotate around U
				addEdge(i, j+1, i+1, j, p)
				addEdge(i+1, j-1, i, j, p) // rotate around D
				addEdge(i+1, j+1, i, j, p)
			}
		}
	}

	if dotCount < 2 {
		fmt.Fprintln(out, -1)
		return
	}

	const INF int64 = 1 << 60
	dist := make([]int64, n*m)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				id := idx(i, j)
				dist[id] = 0
				heap.Push(pq, Item{node: id, dist: 0})
			}
		}
	}

	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		v := it.node
		for _, e := range adj[v] {
			nd := it.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{node: e.to, dist: nd})
			}
		}
	}

	ans := INF
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			id := idx(i, j)
			if j+1 < m {
				id2 := idx(i, j+1)
				sum := dist[id] + dist[id2]
				if sum < ans {
					ans = sum
				}
			}
			if i+1 < n {
				id2 := idx(i+1, j)
				sum := dist[id] + dist[id2]
				if sum < ans {
					ans = sum
				}
			}
		}
	}
	if ans >= INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, ans)
	}
}
