package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 150005
const logN = 18

var (
	adj        [maxN][]edge
	up         [logN][maxN]int
	parentEdge [maxN]int
	depth      [maxN]int
)

type edge struct{ to, id int }

func dfs(v, p int) {
	for _, e := range adj[v] {
		if e.to == p {
			continue
		}
		up[0][e.to] = v
		parentEdge[e.to] = e.id
		depth[e.to] = depth[v] + 1
		for i := 1; i < logN; i++ {
			up[i][e.to] = up[i-1][up[i-1][e.to]]
		}
		dfs(e.to, v)
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := 0; i < logN; i++ {
		if diff&(1<<i) > 0 {
			a = up[i][a]
		}
	}
	if a == b {
		return a
	}
	for i := logN - 1; i >= 0; i-- {
		if up[i][a] != up[i][b] {
			a = up[i][a]
			b = up[i][b]
		}
	}
	return up[0][a]
}

func getPathEdges(u, v int) []int {
	p := lca(u, v)
	res := make([]int, 0)
	x := u
	for x != p {
		res = append(res, parentEdge[x])
		x = up[0][x]
	}
	tmp := make([]int, 0)
	x = v
	for x != p {
		tmp = append(tmp, parentEdge[x])
		x = up[0][x]
	}
	for i := len(tmp) - 1; i >= 0; i-- {
		res = append(res, tmp[i])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	for i := 1; i <= n; i++ {
		adj[i] = nil
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], edge{v, i})
		adj[v] = append(adj[v], edge{u, i})
	}
	depth[1] = 0
	dfs(1, 0)
	edgeTrav := make([][]int, n-1)
	for i := 0; i < m; i++ {
		var s, t int
		fmt.Fscan(reader, &s, &t)
		path := getPathEdges(s, t)
		for _, e := range path {
			edgeTrav[e] = append(edgeTrav[e], i)
		}
	}
	pairCount := make(map[uint64]int)
	ans := 0
	for _, list := range edgeTrav {
		L := len(list)
		for i := 0; i < L; i++ {
			for j := i + 1; j < L; j++ {
				a := list[i]
				b := list[j]
				if a > b {
					a, b = b, a
				}
				key := uint64(a)<<32 | uint64(b)
				pairCount[key]++
				if pairCount[key] == k {
					ans++
				}
			}
		}
	}
	fmt.Println(ans)
}
