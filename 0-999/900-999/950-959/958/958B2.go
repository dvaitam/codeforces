package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	// compute distance from nearest leaf using BFS
	depth := make([]int, n)
	for i := range depth {
		depth[i] = -1
	}
	q := make([]int, 0)
	for i := 0; i < n; i++ {
		if len(g[i]) <= 1 {
			depth[i] = 0
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if depth[to] == -1 {
				depth[to] = depth[v] + 1
				q = append(q, to)
			}
		}
	}

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if depth[order[i]] == depth[order[j]] {
			return order[i] < order[j]
		}
		return depth[order[i]] > depth[order[j]]
	})

	parent := make([]int, n)
	size := make([]int, n)
	leaves := make([]int, n)
	active := make([]bool, n)
	degInside := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
		leaves[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	ans := make([]int, n+1)
	for _, v := range order {
		active[v] = true
		for _, to := range g[v] {
			if active[to] {
				if degInside[v] == 1 {
					rv := find(v)
					leaves[rv]--
				}
				degInside[v]++
				if degInside[to] == 1 {
					rt := find(to)
					leaves[rt]--
				}
				degInside[to]++
				rv := find(v)
				rt := find(to)
				if rv != rt {
					if size[rv] < size[rt] {
						rv, rt = rt, rv
					}
					parent[rt] = rv
					size[rv] += size[rt]
					leaves[rv] += leaves[rt]
				}
			}
		}
		r := find(v)
		if size[r] > ans[leaves[r]] {
			ans[leaves[r]] = size[r]
		}
	}

	for k := 1; k <= n; k++ {
		if ans[k] < ans[k-1] {
			ans[k] = ans[k-1]
		}
	}
	for k := 1; k <= n; k++ {
		if k > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[k])
	}
	fmt.Fprintln(out)
}
