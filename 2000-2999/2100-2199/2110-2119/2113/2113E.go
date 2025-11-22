package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

const maxLog = 18 // enough for n up to 1e5

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, x, y int
		fmt.Fscan(in, &n, &m, &x, &y)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		// binary lifting
		up := make([][]int, maxLog)
		for i := range up {
			up[i] = make([]int, n+1)
		}
		depth := make([]int, n+1)
		// iterative DFS
		stack := []int{1}
		parent := make([]int, n+1)
		parent[1] = 1
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				depth[to] = depth[v] + 1
				up[0][to] = v
				stack = append(stack, to)
			}
		}
		for k := 1; k < maxLog; k++ {
			for v := 1; v <= n; v++ {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}

		lca := func(a, b int) int {
			if depth[a] < depth[b] {
				a, b = b, a
			}
			diff := depth[a] - depth[b]
			for k := 0; diff > 0; k, diff = k+1, diff>>1 {
				if diff&1 == 1 {
					a = up[k][a]
				}
			}
			if a == b {
				return a
			}
			for k := maxLog - 1; k >= 0; k-- {
				if up[k][a] != up[k][b] {
					a = up[k][a]
					b = up[k][b]
				}
			}
			return up[0][a]
		}

		blocked := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			l := lca(a, b)
			path := make([]int, 0)
			u := a
			for u != l {
				path = append(path, u)
				u = up[0][u]
			}
			path = append(path, l)
			temp := make([]int, 0)
			v := b
			for v != l {
				temp = append(temp, v)
				v = up[0][v]
			}
			for i := len(temp) - 1; i >= 0; i-- {
				path = append(path, temp[i])
			}
			for t, node := range path {
				blocked[node] = append(blocked[node], t+1)
			}
		}

		for v := 1; v <= n; v++ {
			if len(blocked[v]) > 1 {
				sort.Ints(blocked[v])
			}
		}

		isBlocked := func(v, t int) bool {
			lst := blocked[v]
			i := sort.SearchInts(lst, t)
			return i < len(lst) && lst[i] == t
		}

		const inf = int(1e18)
		best := make([][]int, n+1)
		for v := 1; v <= n; v++ {
			best[v] = make([]int, len(blocked[v])+1)
			for i := range best[v] {
				best[v][i] = inf
			}
		}

		if isBlocked(x, 1) {
			fmt.Fprintln(out, -1)
			continue
		}
		idxStart := sort.SearchInts(blocked[x], 1)
		best[x][idxStart] = 1
		pq := &MinHeap{}
		heap.Push(pq, Item{time: 1, v: x, idx: idxStart})

		ans := -1
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			t, v, idx := it.time, it.v, it.idx
			if t != best[v][idx] {
				continue
			}
			if v == y {
				ans = t
				break
			}
			nt := t + 1
			// wait
			if !isBlocked(v, nt) {
				idx2 := sort.SearchInts(blocked[v], nt)
				if nt < best[v][idx2] {
					best[v][idx2] = nt
					heap.Push(pq, Item{time: nt, v: v, idx: idx2})
				}
			}
			// move
			for _, to := range adj[v] {
				if isBlocked(to, nt) {
					continue
				}
				idx2 := sort.SearchInts(blocked[to], nt)
				if nt < best[to][idx2] {
					best[to][idx2] = nt
					heap.Push(pq, Item{time: nt, v: to, idx: idx2})
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}

type Item struct {
	time int
	v    int
	idx  int
}

type MinHeap []Item

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].time < h[j].time }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
