package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, k int
	adj  [][]int
	col  []int
	bo   []bool
	la   []int
	dep  []int
	mxd  []int
	path []int
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getcol2(u, p int) {
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		col[v] = 3 - col[u]
		getcol2(v, u)
	}
}

func dfs(u, p int) {
	dep[u] = dep[p] + 1
	mxd[u] = dep[u]
	for _, v := range adj[u] {
		if v == p || bo[v] {
			continue
		}
		dfs(v, u)
		if mxd[v] > mxd[u] {
			mxd[u] = mxd[v]
		}
	}
}

func dfs1(u, p int) {
	if col[p] == k {
		col[u] = 1
	} else {
		col[u] = col[p] + 1
	}
	for _, v := range adj[u] {
		if v == p || bo[v] {
			continue
		}
		dfs1(v, u)
	}
}

func dfs2(u, p int) {
	if col[p] == 1 {
		col[u] = k
	} else {
		col[u] = col[p] - 1
	}
	for _, v := range adj[u] {
		if v == p || bo[v] {
			continue
		}
		dfs2(v, u)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n, &k)
	adj = make([][]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	col = make([]int, n+1)
	bo = make([]bool, n+1)
	la = make([]int, n+1)
	dep = make([]int, n+1)
	mxd = make([]int, n+1)

	if k == 2 {
		col[1] = 1
		getcol2(1, 0)
		fmt.Fprintln(out, "Yes")
		for i := 1; i <= n; i++ {
			fmt.Fprint(out, col[i], " ")
		}
		fmt.Fprintln(out)
		return
	}
	// first BFS from 1
	q := make([]int, 0, n)
	q = append(q, 1)
	bo[1] = true
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if !bo[v] {
				bo[v] = true
				q = append(q, v)
			}
		}
	}
	start := q[len(q)-1]
	// BFS from start
	for i := range bo {
		bo[i] = false
	}
	q = q[:0]
	q = append(q, start)
	bo[start] = true
	la[start] = 0
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if !bo[v] {
				bo[v] = true
				la[v] = u
				q = append(q, v)
			}
		}
	}
	end := q[len(q)-1]
	// build path
	path = make([]int, 0, n)
	for u := end; u != 0; u = la[u] {
		path = append(path, u)
	}
	L := len(path)
	// mark path and initial colors
	for i, u := range path {
		bo[u] = true
		col[u] = ((i + 1) % k) + 1
	}
	// process each path node
	for i, u := range path {
		mxd[u] = 0
		for _, v := range adj[u] {
			if bo[v] {
				continue
			}
			dfs(v, u)
			if mxd[v] > mxd[u] {
				mxd[u] = mxd[v]
			}
		}
		// check impossible
		if mxd[u] > 0 && mxd[u]+(i+1) >= k && mxd[u]+(L-i) >= k {
			fmt.Fprintln(out, "No")
			return
		}
		// color subtrees
		for _, v := range adj[u] {
			if bo[v] {
				continue
			}
			if mxd[u]+(i+1) >= k {
				dfs1(v, u)
			} else {
				dfs2(v, u)
			}
		}
	}
	fmt.Fprintln(out, "Yes")
	for i := 1; i <= n; i++ {
		fmt.Fprint(out, col[i], " ")
	}
	fmt.Fprintln(out)
}
