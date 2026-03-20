package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	cap  int
	cost int
	rev  int
}

type Item struct {
	v    int
	dist int
}

type PQ []Item

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var last_val [100005]int
	var last_mod [7]int
	for i := 0; i <= 100000; i++ {
		last_val[i] = -1
	}
	for i := 0; i < 7; i++ {
		last_mod[i] = -1
	}

	next_val := make([]int, n+1)
	next_mod := make([]int, n+1)
	next_plus := make([]int, n+1)
	next_minus := make([]int, n+1)

	for i := n; i >= 1; i-- {
		val := a[i]
		next_val[i] = last_val[val]
		if val+1 <= 100000 {
			next_plus[i] = last_val[val+1]
		} else {
			next_plus[i] = -1
		}
		if val-1 >= 1 {
			next_minus[i] = last_val[val-1]
		} else {
			next_minus[i] = -1
		}
		next_mod[i] = last_mod[val%7]

		last_val[val] = i
		last_mod[val%7] = i
	}

	numNodes := 3 + 4*n
	S := 0
	S_prime := 1
	T := 2

	adj := make([][]Edge, numNodes)

	addEdge := func(u, v, cap, cost int) {
		adj[u] = append(adj[u], Edge{v, cap, cost, len(adj[v])})
		adj[v] = append(adj[v], Edge{u, 0, -cost, len(adj[u]) - 1})
	}

	addEdge(S, S_prime, 4, 0)

	valBusNode := func(i int) int { return 3 + (i-1)*4 + 0 }
	modBusNode := func(i int) int { return 3 + (i-1)*4 + 1 }
	inNode := func(i int) int { return 3 + (i-1)*4 + 2 }
	outNode := func(i int) int { return 3 + (i-1)*4 + 3 }

	for i := 1; i <= n; i++ {
		addEdge(S_prime, inNode(i), 1, 0)
		addEdge(outNode(i), T, 1, 0)
		addEdge(inNode(i), outNode(i), 1, -1)

		addEdge(valBusNode(i), inNode(i), 1, 0)
		addEdge(modBusNode(i), inNode(i), 1, 0)

		if next_val[i] != -1 {
			addEdge(valBusNode(i), valBusNode(next_val[i]), 4, 0)
		}
		if next_mod[i] != -1 {
			addEdge(modBusNode(i), modBusNode(next_mod[i]), 4, 0)
		}

		if next_plus[i] != -1 {
			addEdge(outNode(i), valBusNode(next_plus[i]), 1, 0)
		}
		if next_minus[i] != -1 {
			addEdge(outNode(i), valBusNode(next_minus[i]), 1, 0)
		}
		if next_mod[i] != -1 {
			addEdge(outNode(i), modBusNode(next_mod[i]), 1, 0)
		}
	}

	h := make([]int, numNodes)
	for i := 0; i < numNodes; i++ {
		h[i] = 1e9
	}
	h[S] = 0

	topo := make([]int, 0, numNodes)
	topo = append(topo, S, S_prime)
	for i := 1; i <= n; i++ {
		topo = append(topo, valBusNode(i), modBusNode(i), inNode(i), outNode(i))
	}
	topo = append(topo, T)

	for _, u := range topo {
		if h[u] != 1e9 {
			for _, e := range adj[u] {
				if e.cap > 0 {
					if h[e.to] > h[u]+e.cost {
						h[e.to] = h[u] + e.cost
					}
				}
			}
		}
	}

	total_cost := 0
	for step := 0; step < 4; step++ {
		dist := make([]int, numNodes)
		for i := 0; i < numNodes; i++ {
			dist[i] = 1e9
		}
		dist[S] = 0

		parentEdge := make([]int, numNodes)
		parentVertex := make([]int, numNodes)
		for i := 0; i < numNodes; i++ {
			parentVertex[i] = -1
		}

		pq := make(PQ, 0)
		heap.Push(&pq, Item{v: S, dist: 0})

		for len(pq) > 0 {
			curr := heap.Pop(&pq).(Item)
			u := curr.v
			d := curr.dist
			if d > dist[u] {
				continue
			}

			for i, e := range adj[u] {
				if e.cap > 0 {
					w := e.cost + h[u] - h[e.to]
					if dist[e.to] > dist[u]+w {
						dist[e.to] = dist[u] + w
						parentVertex[e.to] = u
						parentEdge[e.to] = i
						heap.Push(&pq, Item{v: e.to, dist: dist[e.to]})
					}
				}
			}
		}

		if dist[T] == 1e9 {
			break
		}

		total_cost += dist[T] + h[T]

		curr := T
		for curr != S {
			p := parentVertex[curr]
			idx := parentEdge[curr]

			adj[p][idx].cap -= 1
			revIdx := adj[p][idx].rev
			adj[curr][revIdx].cap += 1

			curr = p
		}

		for i := 0; i < numNodes; i++ {
			if dist[i] != 1e9 {
				h[i] += dist[i]
			}
		}
	}

	fmt.Println(-total_cost)
}