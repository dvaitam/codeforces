package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// encodeRooted returns a canonical string representation of the subtree
// rooted at node with given parent using AHU algorithm.
func encodeRooted(node, parent int, E [][]int) string {
	var children []string
	for _, v := range E[node] {
		if v == parent {
			continue
		}
		children = append(children, encodeRooted(v, node, E))
	}
	sort.Strings(children)
	res := "("
	for _, ch := range children {
		res += ch
	}
	res += ")"
	return res
}

// farthest returns the farthest node from start and parent array.
func farthest(start int, E [][]int) (int, []int) {
	n := len(E)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	q := []int{start}
	parent[start] = -1
	var last int
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		last = x
		for _, y := range E[x] {
			if parent[y] != -2 {
				continue
			}
			parent[y] = x
			q = append(q, y)
		}
	}
	return last, parent
}

// findCenters returns one or two centers of the tree represented by E.
func findCenters(E [][]int) []int {
	if len(E) == 1 {
		return []int{0}
	}
	a, _ := farthest(0, E)
	b, parent := farthest(a, E)
	path := []int{}
	for cur := b; cur != -1; cur = parent[cur] {
		path = append(path, cur)
	}
	l := len(path)
	if l%2 == 1 {
		return []int{path[l/2]}
	}
	return []int{path[l/2], path[l/2-1]}
}

// canonicalForm computes an isomorphism-invariant string for the tree E.
func canonicalForm(E [][]int) string {
	centers := findCenters(E)
	best := ""
	for _, c := range centers {
		code := encodeRooted(c, -1, E)
		if best == "" || code < best {
			best = code
		}
	}
	return best
}

// buildComponent builds the adjacency list of the component reachable from start
// without crossing the banned node. It returns local adjacency and the mapping
// of local index to global node.
func buildComponent(start, banned int, adj [][]int) ([][]int, []int) {
	queue := []int{start}
	visited := map[int]int{start: 0}
	nodes := []int{start}
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range adj[x] {
			if y == banned {
				continue
			}
			if _, ok := visited[y]; !ok {
				visited[y] = len(nodes)
				nodes = append(nodes, y)
				queue = append(queue, y)
			}
		}
	}
	m := len(nodes)
	local := make([][]int, m)
	for i, x := range nodes {
		for _, y := range adj[x] {
			if y == banned {
				continue
			}
			if j, ok := visited[y]; ok {
				local[i] = append(local[i], j)
			}
		}
	}
	return local, nodes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	maxComp := -1
	for u := 0; u < n; u++ {
		if len(adj[u]) < 2 {
			continue
		}
		forms := make([]string, len(adj[u]))
		allSame := true
		for idx, v := range adj[u] {
			localAdj, _ := buildComponent(v, u, adj)
			forms[idx] = canonicalForm(localAdj)
			if idx > 0 && forms[idx] != forms[0] {
				allSame = false
				break
			}
		}
		if allSame && len(forms) > maxComp {
			maxComp = len(forms)
		}
	}

	if maxComp == -1 {
		fmt.Println(-1)
	} else {
		fmt.Println(maxComp)
	}
}
