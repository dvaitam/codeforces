package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const LOG = 17

type Edge struct {
	to int
	w  int
}

var (
	n     int
	adj   [][]int
	up    [][]int
	depth []int
	tin   []int
	timer int
)

func dfs(u, p int) {
	timer++
	tin[u] = timer
	up[u][0] = p
	for i := 1; i < LOG; i++ {
		up[u][i] = up[up[u][i-1]][i-1]
	}
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs(v, u)
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			a = up[a][i]
		}
	}
	if a == b {
		return a
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[a][i] != up[b][i] {
			a = up[a][i]
			b = up[b][i]
		}
	}
	return up[a][0]
}

func buildVirtualTree(nodes []int, imp []bool) (int, map[int][]Edge, bool) {
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	m := len(nodes)
	// add lcas
	extra := make([]int, 0, m*2)
	extra = append(extra, nodes...)
	for i := 0; i < m-1; i++ {
		l := lca(nodes[i], nodes[i+1])
		extra = append(extra, l)
	}
	sort.Slice(extra, func(i, j int) bool { return tin[extra[i]] < tin[extra[j]] })
	uniq := extra[:1]
	for i := 1; i < len(extra); i++ {
		if extra[i] != extra[i-1] {
			uniq = append(uniq, extra[i])
		}
	}
	nodes = uniq

	vt := make(map[int][]Edge, len(nodes))
	stack := []int{}
	impossible := false

	var addEdge func(u, v int)
	addEdge = func(u, v int) {
		w := depth[v] - depth[u]
		if w == 1 && imp[u] && imp[v] {
			impossible = true
		}
		vt[u] = append(vt[u], Edge{v, w})
	}

	for _, v := range nodes {
		if len(stack) == 0 {
			stack = append(stack, v)
			continue
		}
		l := lca(v, stack[len(stack)-1])
		for len(stack) >= 2 && depth[stack[len(stack)-2]] >= depth[l] {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			addEdge(stack[len(stack)-1], u)
		}
		if stack[len(stack)-1] != l {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 || stack[len(stack)-1] != l {
				stack = append(stack, l)
			}
			addEdge(l, u)
		}
		stack = append(stack, v)
	}
	for len(stack) > 1 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		addEdge(stack[len(stack)-1], u)
	}
	root := stack[0]
	return root, vt, impossible
}

func solve(v int, vt map[int][]Edge, imp []bool) (int, bool) {
	cost := 0
	openCnt := 0
	for _, e := range vt[v] {
		c, open := solve(e.to, vt, imp)
		cost += c
		if open {
			if e.w > 1 {
				cost++
			} else {
				if imp[v] {
					cost++
				} else {
					openCnt++
				}
			}
		}
	}
	if imp[v] {
		return cost, true
	}
	if openCnt >= 2 {
		cost++
		return cost, false
	}
	if openCnt == 1 {
		return cost, true
	}
	return cost, false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	up = make([][]int, n+1)
	for i := range up {
		up[i] = make([]int, LOG)
	}
	depth = make([]int, n+1)
	tin = make([]int, n+1)
	timer = 0
	dfs(1, 1)

	var q int
	fmt.Fscan(in, &q)
	imp := make([]bool, n+1)
	nodes := make([]int, 0, n)
	for ; q > 0; q-- {
		var k int
		fmt.Fscan(in, &k)
		nodes = nodes[:0]
		for i := 0; i < k; i++ {
			var x int
			fmt.Fscan(in, &x)
			nodes = append(nodes, x)
			imp[x] = true
		}
		root, vt, bad := buildVirtualTree(nodes, imp)
		if bad {
			fmt.Fprintln(out, -1)
		} else {
			res, open := solve(root, vt, imp)
			if open && !imp[root] {
				res++
			}
			fmt.Fprintln(out, res)
		}
		for _, x := range nodes {
			imp[x] = false
		}
		for v := range vt {
			vt[v] = nil
		}
	}
}
