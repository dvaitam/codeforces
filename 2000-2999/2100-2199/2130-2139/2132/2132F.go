package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

const (
	mod1 = 1000000007
	mod2 = 1000000009
	inf  = int(1e9)
)

type Edge struct {
	to int
	id int
}

type Pair struct {
	a int64
	b int64
}

func addPair(x, y Pair) Pair {
	return Pair{(x.a + y.a) % mod1, (x.b + y.b) % mod2}
}

func mulPair(x, y Pair) Pair {
	return Pair{(x.a * y.a) % mod1, (x.b * y.b) % mod2}
}

func bfs(start int, adj [][]Edge) []int {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	queue := make([]int, 0)
	queue = append(queue, start)
	dist[start] = 0
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, e := range adj[u] {
			v := e.to
			if dist[v] == inf {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}
	return dist
}

type Node struct {
	dist int
	idx  int
	v    int
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].dist == pq[j].dist {
		return pq[i].idx < pq[j].idx
	}
	return pq[i].dist < pq[j].dist
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Node))
}
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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		adj := make([][]Edge, n+1)
		edges := make([][3]int, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			edges[i] = [3]int{u, v, i + 1}
			adj[u] = append(adj[u], Edge{v, i + 1})
			adj[v] = append(adj[v], Edge{u, i + 1})
		}
		dist1 := bfs(1, adj)
		distn := bfs(n, adj)
		D := dist1[n]

		onPath := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			if dist1[i]+distn[i] == D {
				onPath[i] = true
			}
		}

		nodes := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			nodes = append(nodes, i)
		}
		// sort by dist1
		sort.Slice(nodes, func(i, j int) bool {
			if dist1[nodes[i]] == dist1[nodes[j]] {
				return nodes[i] < nodes[j]
			}
			return dist1[nodes[i]] < dist1[nodes[j]]
		})

		dp1 := make([]Pair, n+1)
		dp2 := make([]Pair, n+1)
		dp1[1] = Pair{1, 1}
		for _, u := range nodes {
			if !onPath[u] {
				continue
			}
			for _, e := range adj[u] {
				v := e.to
				if !onPath[v] {
					continue
				}
				if dist1[v] == dist1[u]+1 {
					dp1[v] = addPair(dp1[v], dp1[u])
				}
			}
		}
		dp2[n] = Pair{1, 1}
		sort.Slice(nodes, func(i, j int) bool {
			if dist1[nodes[i]] == dist1[nodes[j]] {
				return nodes[i] < nodes[j]
			}
			return dist1[nodes[i]] > dist1[nodes[j]]
		})
		for _, u := range nodes {
			if !onPath[u] {
				continue
			}
			for _, e := range adj[u] {
				v := e.to
				if !onPath[v] {
					continue
				}
				if dist1[v]+1 == dist1[u] {
					dp2[v] = addPair(dp2[v], dp2[u])
				}
			}
		}
		total := dp1[n]
		critical := make([][3]int, 0)
		for _, e := range edges {
			u, v, id := e[0], e[1], e[2]
			if !onPath[u] || !onPath[v] {
				continue
			}
			if abs(dist1[u]-dist1[v]) != 1 {
				continue
			}
			var from, to int
			if dist1[u] < dist1[v] {
				from, to = u, v
			} else {
				from, to = v, u
			}
			if dp1[from].a == 0 || dp1[from].b == 0 {
				continue
			}
			if dp2[to].a == 0 || dp2[to].b == 0 {
				continue
			}
			prod := mulPair(dp1[from], dp2[to])
			if prod == total {
				critical = append(critical, [3]int{e[0], e[1], id})
			}
		}

		distEdge := make([]int, n+1)
		bestEdge := make([]int, n+1)
		for i := 1; i <= n; i++ {
			distEdge[i] = inf
			bestEdge[i] = inf
		}
		if len(critical) > 0 {
			pq := &PriorityQueue{}
			heap.Init(pq)
			for _, e := range critical {
				u, v, id := e[0], e[1], e[2]
				heap.Push(pq, Node{0, id, u})
				heap.Push(pq, Node{0, id, v})
			}
			for pq.Len() > 0 {
				cur := heap.Pop(pq).(Node)
				if cur.dist > distEdge[cur.v] {
					continue
				}
				if cur.dist == distEdge[cur.v] && cur.idx >= bestEdge[cur.v] {
					continue
				}
				distEdge[cur.v] = cur.dist
				bestEdge[cur.v] = cur.idx
				for _, e := range adj[cur.v] {
					to := e.to
					nd := cur.dist + 1
					if nd < distEdge[to] || (nd == distEdge[to] && cur.idx < bestEdge[to]) {
						heap.Push(pq, Node{nd, cur.idx, to})
					}
				}
			}
		}

		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var c int
			fmt.Fscan(in, &c)
			if distEdge[c] == inf {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, bestEdge[c])
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
