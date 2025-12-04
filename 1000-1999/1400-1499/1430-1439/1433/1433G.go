package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF = 1000000000 // A large enough value to represent infinity

type Edge struct {
	to, weight int
}

type Item struct {
	node, priority int
	index          int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest priority (smallest distance) item as first.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func dijkstra(n, start int, adj [][]Edge) []int {
	d := make([]int, n) // Changed to n, for 0-based indexing
	for i := 0; i < n; i++ { // Changed loop to 0 to n-1
		d[i] = INF
	}
	d[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{node: start, priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		u := item.node

		if item.priority > d[u] {
			continue
		}

		for _, e := range adj[u] {
			if d[u]+e.weight < d[e.to] {
				d[e.to] = d[u] + e.weight
				heap.Push(&pq, &Item{node: e.to, priority: d[e.to]})
			}
		}
	}
	return d
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)

	adj := make([][]Edge, n) // Changed to n, for 0-based indexing
	type Road struct {
		u, v, w int
	}
	roads := make([]Road, m)

	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		u-- // Convert to 0-based
		v-- // Convert to 0-based
		adj[u] = append(adj[u], Edge{to: v, weight: w})
		adj[v] = append(adj[v], Edge{to: u, weight: w})
		roads[i] = Road{u, v, w}
	}

	routes := make([][2]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &routes[i][0], &routes[i][1])
		routes[i][0]-- // Convert to 0-based
		routes[i][1]-- // Convert to 0-based
	}

	dist := make([][]int, n) // Changed to n, for 0-based indexing
	for i := 0; i < n; i++ { // Changed loop to 0 to n-1
		dist[i] = dijkstra(n, i, adj)
	}

	ans := 0
	for i := 0; i < k; i++ {
		ans += dist[routes[i][0]][routes[i][1]]
	}

	for _, r := range roads {
		cur := 0
		u, v := r.u, r.v
		for i := 0; i < k; i++ {
			a, b := routes[i][0], routes[i][1]
			d1 := dist[a][b]
			d2 := dist[a][u] + dist[v][b]
			d3 := dist[a][v] + dist[u][b]
			mn := d1
			if d2 < mn {
				mn = d2
			}
			if d3 < mn {
				mn = d3
			}
			cur += mn
		}
		if cur < ans {
			ans = cur
		}
	}

	fmt.Fprintln(writer, ans)
}