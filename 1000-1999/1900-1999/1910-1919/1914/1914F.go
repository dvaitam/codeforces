package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Item represents a node in the priority queue sorted by depth (max-heap).
type Item struct {
	depth int
	node  int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].depth > h[j].depth }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		parent := make([]int, n+1)
		g := make([][]int, n+1)
		deg := make([]int, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &parent[i])
			p := parent[i]
			g[p] = append(g[p], i)
			g[i] = append(g[i], p)
			deg[p]++
			deg[i]++
		}

		// Iterative DFS to compute tin, tout and depth.
		tin := make([]int, n+1)
		tout := make([]int, n+1)
		depth := make([]int, n+1)
		type node struct{ v, p, st int }
		stack := []node{{1, 0, 0}}
		timer := 0
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if cur.st == 0 {
				if cur.v == 1 {
					depth[cur.v] = 0
				} else {
					depth[cur.v] = depth[cur.p] + 1
				}
				tin[cur.v] = timer
				timer++
				stack = append(stack, node{cur.v, cur.p, 1})
				for i := len(g[cur.v]) - 1; i >= 0; i-- {
					to := g[cur.v][i]
					if to == cur.p {
						continue
					}
					stack = append(stack, node{to, cur.v, 0})
				}
			} else {
				tout[cur.v] = timer - 1
			}
		}

		used := make([]bool, n+1)
		pq := &MaxHeap{}
		heap.Init(pq)
		for i := 2; i <= n; i++ {
			if deg[i] == 1 {
				heap.Push(pq, Item{depth: depth[i], node: i})
			}
		}
		pairs := 0
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			u := it.node
			if used[u] || deg[u] != 1 {
				continue
			}
			var temp []Item
			cand := 0
			for pq.Len() > 0 {
				it2 := heap.Pop(pq).(Item)
				v := it2.node
				if used[v] || deg[v] != 1 {
					continue
				}
				// Skip if v is an ancestor of u.
				if tin[v] <= tin[u] && tout[u] <= tout[v] {
					temp = append(temp, it2)
					continue
				}
				cand = v
				break
			}
			for _, x := range temp {
				heap.Push(pq, x)
			}
			if cand == 0 {
				continue
			}
			pairs++
			used[u], used[cand] = true, true
			for _, x := range []int{u, cand} {
				for _, nb := range g[x] {
					deg[nb]--
					if nb != 1 && !used[nb] && deg[nb] == 1 {
						heap.Push(pq, Item{depth: depth[nb], node: nb})
					}
				}
			}
		}
		fmt.Fprintln(out, pairs)
	}
}
