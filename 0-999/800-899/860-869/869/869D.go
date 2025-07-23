package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to    int
	extra bool
}

var (
	adj     [][]Edge
	weights []int64
	mod     int64 = 1000000007
	visited []bool
	start   int
	result  int64
)

func dfs(u int, used bool) {
	for _, e := range adj[u] {
		v := e.to
		if visited[v] {
			continue
		}
		visited[v] = true
		nu := used || e.extra
		if nu {
			result = (result + (weights[start]%mod)*(weights[v]%mod)) % mod
		}
		dfs(v, nu)
		visited[v] = false
	}
}

func depth(x int) int {
	d := 0
	for x > 1 {
		x >>= 1
		d++
	}
	return d
}

func lca(a, b int) int {
	for a != b {
		if a > b {
			a >>= 1
		} else {
			b >>= 1
		}
	}
	return a
}

func subtreeSize(n int64, x int64) int64 {
	cnt := int64(0)
	l := x
	r := x
	for l <= n {
		if r > n {
			r = n
		}
		cnt += r - l + 1
		l = l * 2
		r = r*2 + 1
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	var m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	type pair struct{ u, v int }
	edges := make([]pair, m)
	nodesMap := map[int]bool{1: true}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = pair{u, v}
		nodesMap[u] = true
		nodesMap[v] = true
	}
	// close under LCA
	changed := true
	for changed {
		changed = false
		curNodes := make([]int, 0, len(nodesMap))
		for x := range nodesMap {
			curNodes = append(curNodes, x)
		}
		for i := 0; i < len(curNodes); i++ {
			for j := i + 1; j < len(curNodes); j++ {
				l := lca(curNodes[i], curNodes[j])
				if !nodesMap[l] {
					nodesMap[l] = true
					changed = true
				}
			}
		}
	}

	// convert to slice and sort by value (not really needed but stable)
	nodes := make([]int, 0, len(nodesMap))
	for x := range nodesMap {
		nodes = append(nodes, x)
	}
	// simple insertion sort by value
	for i := 1; i < len(nodes); i++ {
		v := nodes[i]
		j := i - 1
		for j >= 0 && nodes[j] > v {
			nodes[j+1] = nodes[j]
			j--
		}
		nodes[j+1] = v
	}

	id := make(map[int]int)
	for i, x := range nodes {
		id[x] = i
	}
	k := len(nodes)
	adj = make([][]Edge, k)
	children := make([][]int, k)

	// tree edges from each node to nearest ancestor in set
	for _, x := range nodes {
		if x == 1 {
			continue
		}
		p := x >> 1
		for !nodesMap[p] {
			p >>= 1
		}
		a := id[x]
		b := id[p]
		adj[a] = append(adj[a], Edge{b, false})
		adj[b] = append(adj[b], Edge{a, false})
		children[b] = append(children[b], a)
	}

	// extra edges
	for _, e := range edges {
		a := id[e.u]
		b := id[e.v]
		adj[a] = append(adj[a], Edge{b, true})
		adj[b] = append(adj[b], Edge{a, true})
	}

	weights = make([]int64, k)
	for i, x := range nodes {
		total := subtreeSize(n, int64(x))
		for _, ch := range children[i] {
			total -= subtreeSize(n, int64(nodes[ch]))
		}
		weights[i] = total % mod
	}

	result = 0
	visited = make([]bool, k)
	for s := 0; s < k; s++ {
		start = s
		visited[s] = true
		dfs(s, false)
		visited[s] = false
	}

	nmod := (n % mod)
	base := (nmod * nmod) % mod
	ans := (base + result%mod) % mod
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
