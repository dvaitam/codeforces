package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
}

var (
	n, m      int
	w         []int64
	g         [][]Edge
	edges     [][2]int
	disc, low []int
	timer     int
	isBridge  []bool

	comp       []int
	compWeight []int64
	compSize   []int
	comps      int

	tree     [][]int
	good     []bool
	dpRet    []int64
	cycleSub []bool
	dpBest   []int64
)

func dfsBridge(u, parentEdge int) {
	timer++
	disc[u] = timer
	low[u] = timer
	for _, e := range g[u] {
		v, id := e.to, e.id
		if id == parentEdge {
			continue
		}
		if disc[v] == 0 {
			dfsBridge(v, id)
			if low[v] < low[u] {
				low[u] = low[v]
			}
			if low[v] > disc[u] {
				isBridge[id] = true
			}
		} else if disc[v] < low[u] {
			low[u] = disc[v]
		}
	}
}

func dfsComp(u int) {
	comp[u] = comps
	compWeight[comps] += w[u]
	compSize[comps]++
	for _, e := range g[u] {
		if isBridge[e.id] {
			continue
		}
		v := e.to
		if comp[v] == 0 {
			dfsComp(v)
		}
	}
}

func dfs1(u, p int) {
	dpRet[u] = compWeight[u]
	cycleSub[u] = good[u]
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		dfs1(v, u)
		if cycleSub[v] {
			dpRet[u] += dpRet[v]
			cycleSub[u] = true
		}
	}
}

func dfs2(u, p int) {
	dpBest[u] = dpRet[u]
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		dfs2(v, u)
		var cand int64
		if cycleSub[v] {
			cand = dpRet[u] - dpRet[v] + dpBest[v]
		} else {
			cand = dpRet[u] + dpBest[v]
		}
		if cand > dpBest[u] {
			dpBest[u] = cand
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &m)
	w = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &w[i])
	}
	g = make([][]Edge, n+1)
	edges = make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		edges[i] = [2]int{u, v}
		g[u] = append(g[u], Edge{v, i})
		g[v] = append(g[v], Edge{u, i})
	}
	var s int
	fmt.Fscan(reader, &s)

	disc = make([]int, n+1)
	low = make([]int, n+1)
	isBridge = make([]bool, m)

	timer = 0
	dfsBridge(1, -1)
	for i := 1; i <= n; i++ {
		if disc[i] == 0 {
			dfsBridge(i, -1)
		}
	}

	comp = make([]int, n+1)
	compWeight = make([]int64, n+1) // temporary bigger size
	compSize = make([]int, n+1)
	comps = 0
	for i := 1; i <= n; i++ {
		if comp[i] == 0 {
			comps++
			dfsComp(i)
		}
	}

	compWeight = compWeight[:comps+1]
	compSize = compSize[:comps+1]
	good = make([]bool, comps+1)
	for i := 1; i <= comps; i++ {
		if compSize[i] > 1 {
			good[i] = true
		}
	}

	tree = make([][]int, comps+1)
	for i, e := range edges {
		_ = i
		u, v := e[0], e[1]
		cu, cv := comp[u], comp[v]
		if cu != cv {
			tree[cu] = append(tree[cu], cv)
			tree[cv] = append(tree[cv], cu)
		}
	}

	dpRet = make([]int64, comps+1)
	dpBest = make([]int64, comps+1)
	cycleSub = make([]bool, comps+1)

	root := comp[s]
	dfs1(root, 0)
	dfs2(root, 0)

	fmt.Println(dpBest[root])
}
