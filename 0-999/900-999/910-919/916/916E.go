package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, q       int
	adj        [][]int
	LOG        int
	up         [][]int
	tin, tout  []int
	depth      []int
	timer      int
	bit1, bit2 []int64
)

func add(bit []int64, idx int, val int64) {
	for idx <= n {
		bit[idx] += val
		idx += idx & -idx
	}
}

func rangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	add(bit1, l, val)
	add(bit1, r+1, -val)
	add(bit2, l, val*int64(l-1))
	add(bit2, r+1, -val*int64(r))
}

func prefix(bit []int64, idx int) int64 {
	var res int64
	for idx > 0 {
		res += bit[idx]
		idx -= idx & -idx
	}
	return res
}

func prefixSum(idx int) int64 {
	return prefix(bit1, idx)*int64(idx) - prefix(bit2, idx)
}

func rangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return prefixSum(r) - prefixSum(l-1)
}

func dfs(v, p int) {
	timer++
	tin[v] = timer
	up[0][v] = p
	for i := 1; i < LOG; i++ {
		if up[i-1][v] != 0 {
			up[i][v] = up[i-1][up[i-1][v]]
		}
	}
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		depth[to] = depth[v] + 1
		dfs(to, v)
	}
	tout[v] = timer
}

func isAncestor(u, v int) bool {
	return tin[u] <= tin[v] && tout[v] <= tout[u]
}

func lca(u, v int) int {
	if isAncestor(u, v) {
		return u
	}
	if isAncestor(v, u) {
		return v
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != 0 && !isAncestor(up[i][u], v) {
			u = up[i][u]
		}
	}
	return up[0][u]
}

func jump(u, v int) int { // v is ancestor of u, v!=u
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != 0 && depth[up[i][u]] > depth[v] {
			u = up[i][u]
		}
	}
	return u
}

func lcaRoot(u, v, r int) int {
	a := lca(u, v)
	b := lca(u, r)
	c := lca(v, r)
	if a == b {
		return c
	}
	if a == c {
		return b
	}
	return a
}

func addSubtree(v int, val int64, root int) {
	if v == root {
		rangeAdd(1, n, val)
	} else if !isAncestor(v, root) {
		rangeAdd(tin[v], tout[v], val)
	} else {
		w := jump(root, v)
		rangeAdd(1, n, val)
		rangeAdd(tin[w], tout[w], -val)
	}
}

func sumSubtree(v int, root int) int64 {
	if v == root {
		return rangeSum(1, n)
	} else if !isAncestor(v, root) {
		return rangeSum(tin[v], tout[v])
	}
	w := jump(root, v)
	return rangeSum(1, n) - rangeSum(tin[w], tout[w])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n, &q)
	adj = make([][]int, n+1)
	values := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &values[i])
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	LOG = 1
	for (1 << LOG) <= n {
		LOG++
	}
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	tin = make([]int, n+1)
	tout = make([]int, n+1)
	depth = make([]int, n+1)
	timer = 0
	dfs(1, 0)

	bit1 = make([]int64, n+2)
	bit2 = make([]int64, n+2)
	for i := 1; i <= n; i++ {
		rangeAdd(tin[i], tin[i], values[i])
	}

	root := 1
	for ; q > 0; q-- {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			fmt.Fscan(in, &root)
		} else if tp == 2 {
			var u, v int
			var x int64
			fmt.Fscan(in, &u, &v, &x)
			l := lcaRoot(u, v, root)
			addSubtree(l, x, root)
		} else if tp == 3 {
			var v int
			fmt.Fscan(in, &v)
			ans := sumSubtree(v, root)
			fmt.Fprintln(out, ans)
		}
	}
}
