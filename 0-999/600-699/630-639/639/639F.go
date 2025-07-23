package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Implementation to solve Codeforces problem 639F - Bear and Chemistry
// based only on the textual statement available in this repository.

// Edge structure for base graph
type Edge struct {
	to int
	id int
}

var (
	n, m, q     int
	g           [][]Edge
	edges       [][2]int
	timer       int
	tin, low    []int
	isBridge    []bool
	comp        []int
	compCnt     int
	tree        [][]int
	tinT, toutT []int
	depthT      []int
	up          [][]int
	rootOf      []int
)

const LOG = 20

func addBaseEdge(u, v, id int) {
	g[u] = append(g[u], Edge{v, id})
	g[v] = append(g[v], Edge{u, id})
}

func dfsBridge(v, pe int) {
	timer++
	tin[v] = timer
	low[v] = timer
	for _, e := range g[v] {
		if e.id == pe {
			continue
		}
		if tin[e.to] == 0 {
			dfsBridge(e.to, e.id)
			if low[e.to] > tin[v] {
				isBridge[e.id] = true
			}
			if low[e.to] < low[v] {
				low[v] = low[e.to]
			}
		} else {
			if tin[e.to] < low[v] {
				low[v] = tin[e.to]
			}
		}
	}
}

func dfsComp(v, c int) {
	comp[v] = c
	for _, e := range g[v] {
		if isBridge[e.id] || comp[e.to] != 0 {
			continue
		}
		dfsComp(e.to, c)
	}
}

func dfsTree(v, p, root int) {
	timer++
	tinT[v] = timer
	up[0][v] = p
	rootOf[v] = root
	for i := 1; i < LOG; i++ {
		if up[i-1][v] != -1 {
			up[i][v] = up[i-1][up[i-1][v]]
		} else {
			up[i][v] = -1
		}
	}
	for _, to := range tree[v] {
		if to == p {
			continue
		}
		depthT[to] = depthT[v] + 1
		dfsTree(to, v, root)
	}
	toutT[v] = timer
}

func isAncestor(u, v int) bool {
	return tinT[u] <= tinT[v] && toutT[v] <= toutT[u]
}

func lca(u, v int) int {
	if rootOf[u] != rootOf[v] {
		return -1
	}
	if isAncestor(u, v) {
		return u
	}
	if isAncestor(v, u) {
		return v
	}
	for i := LOG - 1; i >= 0; i-- {
		pu := up[i][u]
		if pu != -1 && !isAncestor(pu, v) {
			u = pu
		}
	}
	return up[0][u]
}

// decode applies the online transformation based on current R
func decode(x, R int) int {
	x--
	x = (x + R) % n
	return x + 1
}

type Edge2 struct {
	to int
	id int
}

// buildBCC computes 2-edge-connected component ids for a given graph
func buildBCC(adj [][]Edge2) []int {
	k := len(adj)
	tin := make([]int, k)
	low := make([]int, k)
	visited := make([]bool, k)
	totalEdges := 0
	for _, a := range adj {
		totalEdges += len(a)
	}
	totalEdges /= 2
	bridges := make([]bool, totalEdges)
	timer := 0
	var dfs func(int, int)
	dfs = func(v, pe int) {
		timer++
		visited[v] = true
		tin[v] = timer
		low[v] = timer
		for _, e := range adj[v] {
			if e.id == pe {
				continue
			}
			if !visited[e.to] {
				dfs(e.to, e.id)
				if low[e.to] > tin[v] {
					bridges[e.id] = true
				}
				if low[e.to] < low[v] {
					low[v] = low[e.to]
				}
			} else {
				if tin[e.to] < low[v] {
					low[v] = tin[e.to]
				}
			}
		}
	}
	for i := 0; i < k; i++ {
		if !visited[i] {
			dfs(i, -1)
		}
	}
	comp := make([]int, k)
	cid := 0
	var dfs2 func(int)
	dfs2 = func(v int) {
		comp[v] = cid
		for _, e := range adj[v] {
			if comp[e.to] == 0 && !bridges[e.id] {
				dfs2(e.to)
			}
		}
	}
	for i := 0; i < k; i++ {
		if comp[i] == 0 {
			cid++
			dfs2(i)
		}
	}
	return comp
}

// buildVirtualTree constructs a virtual tree for nodes belonging to the same base tree
func buildVirtualTree(nodes []int) ([]int, [][2]int) {
	sort.Slice(nodes, func(i, j int) bool { return tinT[nodes[i]] < tinT[nodes[j]] })
	// add LCAs of consecutive nodes
	lcas := make([]int, 0, len(nodes)-1)
	for i := 0; i+1 < len(nodes); i++ {
		l := lca(nodes[i], nodes[i+1])
		if l != -1 {
			lcas = append(lcas, l)
		}
	}
	nodes = append(nodes, lcas...)
	sort.Slice(nodes, func(i, j int) bool { return tinT[nodes[i]] < tinT[nodes[j]] })
	// unique
	uniq := nodes[:0]
	for _, v := range nodes {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	nodes = uniq
	if len(nodes) == 1 {
		return nodes, nil
	}
	st := []int{nodes[0]}
	edges := make([][2]int, 0)
	for i := 1; i < len(nodes); i++ {
		v := nodes[i]
		l := lca(v, st[len(st)-1])
		if l == -1 {
			// different trees, shouldn't happen here
			st = append(st, v)
			continue
		}
		for len(st) >= 2 && depthT[st[len(st)-2]] >= depthT[l] {
			a := st[len(st)-1]
			st = st[:len(st)-1]
			b := st[len(st)-1]
			edges = append(edges, [2]int{a, b})
		}
		if st[len(st)-1] != l {
			edges = append(edges, [2]int{st[len(st)-1], l})
			st[len(st)-1] = l
		}
		st = append(st, v)
	}
	for len(st) > 1 {
		a := st[len(st)-1]
		st = st[:len(st)-1]
		b := st[len(st)-1]
		edges = append(edges, [2]int{a, b})
	}
	return nodes, edges
}

func checkQuery(favs []int, qEdges [][2]int) bool {
	if len(favs) <= 1 {
		return true
	}
	groupNodes := make(map[int][]int)
	for _, x := range favs {
		groupNodes[rootOf[x]] = append(groupNodes[rootOf[x]], x)
	}
	for _, e := range qEdges {
		groupNodes[rootOf[e[0]]] = append(groupNodes[rootOf[e[0]]], e[0])
		groupNodes[rootOf[e[1]]] = append(groupNodes[rootOf[e[1]]], e[1])
	}
	nodeSet := make(map[int]struct{})
	edges := make([][2]int, 0)
	for _, arr := range groupNodes {
		sort.Slice(arr, func(i, j int) bool { return tinT[arr[i]] < tinT[arr[j]] })
		arr2, e := buildVirtualTree(arr)
		for _, v := range arr2 {
			nodeSet[v] = struct{}{}
		}
		edges = append(edges, e...)
	}
	for _, e := range qEdges {
		nodeSet[e[0]] = struct{}{}
		nodeSet[e[1]] = struct{}{}
		edges = append(edges, e)
	}
	nodeList := make([]int, 0, len(nodeSet))
	for v := range nodeSet {
		nodeList = append(nodeList, v)
	}
	sort.Ints(nodeList)
	idx := make(map[int]int, len(nodeList))
	for i, v := range nodeList {
		idx[v] = i
	}
	adj := make([][]Edge2, len(nodeList))
	for id, e := range edges {
		u := idx[e[0]]
		v := idx[e[1]]
		adj[u] = append(adj[u], Edge2{v, id})
		adj[v] = append(adj[v], Edge2{u, id})
	}
	compArr := buildBCC(adj)
	first := compArr[idx[favs[0]]]
	for _, x := range favs[1:] {
		if compArr[idx[x]] != first {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m, &q)
	g = make([][]Edge, n)
	edges = make([][2]int, m)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		edges[i] = [2]int{a, b}
		addBaseEdge(a, b, i)
	}
	tin = make([]int, n)
	low = make([]int, n)
	isBridge = make([]bool, m)
	timer = 0
	for i := 0; i < n; i++ {
		if tin[i] == 0 {
			dfsBridge(i, -1)
		}
	}
	comp = make([]int, n)
	compCnt = 0
	for i := 0; i < n; i++ {
		if comp[i] == 0 {
			compCnt++
			dfsComp(i, compCnt)
		}
	}
	tree = make([][]int, compCnt)
	for i, e := range edges {
		if isBridge[i] {
			u := comp[e[0]] - 1
			v := comp[e[1]] - 1
			tree[u] = append(tree[u], v)
			tree[v] = append(tree[v], u)
		}
	}
	tinT = make([]int, compCnt)
	toutT = make([]int, compCnt)
	depthT = make([]int, compCnt)
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, compCnt)
		for j := 0; j < compCnt; j++ {
			up[i][j] = -1
		}
	}
	rootOf = make([]int, compCnt)
	timer = 0
	for i := 0; i < compCnt; i++ {
		if tinT[i] == 0 {
			dfsTree(i, -1, i)
		}
	}

	R := 0
	for qi := 1; qi <= q; qi++ {
		var ni, mi int
		if _, err := fmt.Fscan(reader, &ni, &mi); err != nil {
			return
		}
		favs := make([]int, ni)
		for i := 0; i < ni; i++ {
			var x int
			fmt.Fscan(reader, &x)
			x = decode(x, R)
			x--
			favs[i] = comp[x] - 1
		}
		qEdges := make([][2]int, mi)
		for i := 0; i < mi; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			a = decode(a, R)
			b = decode(b, R)
			a--
			b--
			qEdges[i] = [2]int{comp[a] - 1, comp[b] - 1}
		}
		if checkQuery(favs, qEdges) {
			fmt.Fprintln(writer, "YES")
			R += qi
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
