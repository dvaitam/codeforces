package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Pair represents an interval endpoint or node identifier
type Pair struct {
	first, second int
}

// Edge connects two nodes in the tree
type Edge struct {
	u, v Pair
}

var diag []Pair
var edges []Edge

// dfs builds the tree from nested intervals
func dfs(pIdx *int) Pair {
	c := *pIdx
	*pIdx++
	var neig []Pair
	b, e := 0, 0
	// collect children intervals
	for *pIdx < len(diag) && diag[*pIdx].first < diag[c].second {
		b = diag[*pIdx].first
		e = diag[*pIdx].second
		child := dfs(pIdx)
		neig = append(neig, child)
	}
	// current node id
	id := Pair{diag[c].second, diag[c].second - 1}
	if e == diag[c].second {
		id.second = b
	}
	// add edges between id and each child
	for _, nn := range neig {
		edges = append(edges, Edge{u: id, v: nn})
	}
	return id
}

// centroid decomposition variables
var (
	mk     []bool
	q      []int
	parent []int
	sz     []int
	mc     []int
	answer []int
)

// centroid finds centroid of subtree rooted at c
func centroid(c int, adj [][]int) int {
	b, e := 0, 0
	q[e] = c
	parent[c] = -1
	sz[c] = 1
	mc[c] = 0
	e++
	for b < e {
		u := q[b]
		b++
		for _, v := range adj[u] {
			if v != parent[u] && !mk[v] {
				parent[v] = u
				sz[v] = 1
				mc[v] = 0
				q[e] = v
				e++
			}
		}
	}
	for i := e - 1; i >= 0; i-- {
		u := q[i]
		bc := e - sz[u]
		if mc[u] > bc {
			bc = mc[u]
		}
		if 2*bc <= e {
			return u
		}
		if parent[u] >= 0 {
			sz[parent[u]] += sz[u]
			if sz[u] > mc[parent[u]] {
				mc[parent[u]] = sz[u]
			}
		}
	}
	return -1
}

// calc performs centroid decomposition and labels nodes
func calc(u int, adj [][]int, level int) {
	c := centroid(u, adj)
	mk[c] = true
	answer[c] = level
	for _, v := range adj[c] {
		if !mk[v] {
			calc(v, adj, level+1)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	if m == 0 {
		fmt.Fprintln(writer, 1)
		return
	}
	diag = make([]Pair, m+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		if u > v {
			u, v = v, u
		}
		diag[i] = Pair{u, v}
	}
	// root interval
	diag[m] = Pair{1, n}
	// sort by start asc, end desc
	sort.Slice(diag, func(i, j int) bool {
		if diag[i].first == diag[j].first {
			return diag[i].second > diag[j].second
		}
		return diag[i].first < diag[j].first
	})
	// build edges
	edges = make([]Edge, 0, m)
	pIdx := 0
	dfs(&pIdx)
	// collect unique nodes
	nodes := make([]Pair, 0, len(edges)*2)
	for _, e := range edges {
		nodes = append(nodes, e.u, e.v)
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].first == nodes[j].first {
			return nodes[i].second < nodes[j].second
		}
		return nodes[i].first < nodes[j].first
	})
	uniq := nodes[:0]
	for i, p := range nodes {
		if i == 0 || p != nodes[i-1] {
			uniq = append(uniq, p)
		}
	}
	nodes = uniq
	total := len(nodes)
	// map nodes to indices
	idx := make(map[Pair]int, total)
	for i, p := range nodes {
		idx[p] = i
	}
	// build adjacency list
	adj := make([][]int, total)
	for _, e := range edges {
		u := idx[e.u]
		v := idx[e.v]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// init centroid variables
	mk = make([]bool, total)
	q = make([]int, total)
	parent = make([]int, total)
	sz = make([]int, total)
	mc = make([]int, total)
	answer = make([]int, total)
	// decompose and print
	calc(0, adj, 1)
	for i := 0; i < total; i++ {
		fmt.Fprint(writer, answer[i])
		if i+1 < total {
			writer.WriteByte(' ')
		}
	}
	writer.WriteByte('\n')
}
