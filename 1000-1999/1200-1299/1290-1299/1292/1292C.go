package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program attempts to maximize the sum of mex values for all pairs
// of nodes in a tree as described in problemC.txt. Each edge is assigned
// a distinct integer from 0..n-2. We follow a greedy approach: edges with
// larger pair contributions (u_subtree_size * v_subtree_size) receive
// smaller labels.

// The idea is that the number of pairs whose path contains a given edge is
// size * (n - size). Assigning a small label to an edge with many such pairs
// should increase the overall sum of mex values.

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	type edge struct{ u, v int }
	edges := make([]edge, n-1)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = edge{u, v}
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	// compute sizes for edges by rooting the tree at 0
	parent := make([]int, n)
	order := []int{0}
	parent[0] = -1
	for idx := 0; idx < len(order); idx++ {
		v := order[idx]
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			order = append(order, to)
		}
	}
	size := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			size[v] += size[to]
		}
	}

	edgeSize := make([]int64, n-1)
	for i, e := range edges {
		a, b := e.u, e.v
		if parent[b] == a {
			edgeSize[i] = int64(size[b]) * int64(n-size[b])
		} else if parent[a] == b {
			edgeSize[i] = int64(size[a]) * int64(n-size[a])
		} else {
			// should not happen
		}
	}

	idxs := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i, j int) bool {
		return edgeSize[idxs[i]] > edgeSize[idxs[j]]
	})

	label := make([]int, n-1)
	for rank, id := range idxs {
		label[id] = rank
	}

	// build adjacency with labels
	lg := make([][]struct{ to, val int }, n)
	for i, e := range edges {
		lg[e.u] = append(lg[e.u], struct{ to, val int }{e.v, label[i]})
		lg[e.v] = append(lg[e.v], struct{ to, val int }{e.u, label[i]})
	}

	// function to compute mex for path u,v
	var mexPath func(int, int) int
	mexPath = func(u, v int) int {
		// simple BFS since n <= 3000
		queue := []int{u}
		par := make([]int, n)
		val := make([]int, n)
		for i := 0; i < n; i++ {
			par[i] = -2
		}
		par[u] = -1
		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			if x == v {
				break
			}
			for _, e := range lg[x] {
				if par[e.to] != -2 {
					continue
				}
				par[e.to] = x
				val[e.to] = e.val
				queue = append(queue, e.to)
			}
		}
		used := make(map[int]struct{})
		cur := v
		for par[cur] != -1 {
			used[val[cur]] = struct{}{}
			cur = par[cur]
		}
		m := 0
		for {
			if _, ok := used[m]; !ok {
				return m
			}
			m++
		}
	}

	var ans int64
	for u := 0; u < n; u++ {
		for v := u + 1; v < n; v++ {
			ans += int64(mexPath(u, v))
		}
	}
	fmt.Println(ans)
}
