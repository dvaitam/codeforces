package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

type Edge struct {
	to int
	id int
}

var (
	g     [][]Edge
	edges [][2]int
	paths [][]int
	n, k  int
)

func getPath(a, b int) []int {
	parent := make([]int, n+1)
	pe := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = -1
	}
	q := []int{a}
	parent[a] = a
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			break
		}
		for _, e := range g[v] {
			if parent[e.to] == -1 {
				parent[e.to] = v
				pe[e.to] = e.id
				q = append(q, e.to)
			}
		}
	}
	var res []int
	cur := b
	for cur != a {
		id := pe[cur]
		res = append(res, id)
		if edges[id][0] == cur {
			cur = edges[id][1]
		} else {
			cur = edges[id][0]
		}
	}
	return res
}

func valid(sel [][]int, m int) bool {
	cnt := make([]int, m)
	all := false
	for _, p := range sel {
		for _, id := range p {
			cnt[id]++
		}
	}
	for _, c := range cnt {
		if c != 0 && c != 1 && c != k {
			return false
		}
		if c == k {
			all = true
		}
	}
	return all
}

func dfs(idx int, cur [][]int, m int, ans *int) {
	if idx == k {
		if valid(cur, m) {
			*ans = (*ans + 1) % mod
		}
		return
	}
	for _, p := range paths {
		cur[idx] = p
		dfs(idx+1, cur, m, ans)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &k)
	g = make([][]Edge, n+1)
	edges = make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = [2]int{u, v}
		g[u] = append(g[u], Edge{v, i})
		g[v] = append(g[v], Edge{u, i})
	}
	// enumerate all simple paths by pair of vertices
	for i := 1; i <= n; i++ {
		paths = append(paths, []int{}) // path of length 0
		for j := i + 1; j <= n; j++ {
			p := getPath(i, j)
			paths = append(paths, p)
		}
	}
	m := len(edges)
	cur := make([][]int, k)
	ans := 0
	dfs(0, cur, m, &ans)
	fmt.Println(ans)
}
