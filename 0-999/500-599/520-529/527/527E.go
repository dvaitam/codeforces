package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	// adjacency using head pointers, 1-indexed
	adj := make([]int, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = -1
	}
	// degree count
	apr := make([]int, n+1)
	// storage for odd-degree vertices
	tem := make([]int, 0, n)
	// edge lists and flags
	to := make([]int, 0, (m+n)*4)
	nxt := make([]int, 0, (m+n)*4)
	vis := make([]bool, 0, (m+n)*4)
	vise := make([]bool, 0, (m+n)*4)
	// undirected edge count
	cntege := 0
	// helper to add a directed edge
	addDir := func(u, v int) {
		e := len(to)
		to = append(to, v)
		nxt = append(nxt, adj[u])
		adj[u] = e
		vis = append(vis, false)
		vise = append(vise, false)
	}
	// add undirected edge
	addEdge := func(u, v int) {
		addDir(u, v)
		addDir(v, u)
		cntege++
	}
	// read input edges
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		addEdge(u, v)
		apr[u]++
		apr[v]++
	}
	// collect odd-degree vertices
	for i := 1; i <= n; i++ {
		if apr[i]%2 != 0 {
			tem = append(tem, i)
		}
	}
	// pair odd-degree vertices
	lt := len(tem)
	for i := 0; i < lt; i += 2 {
		if i+1 < lt {
			addEdge(tem[i], tem[i+1])
		} else {
			// self-loop pairing
			addEdge(tem[i], tem[i])
		}
	}
	// prepare iterator for Euler traversal
	iter := make([]int, n+1)
	copy(iter, adj)
	// stacks for nodes and edges
	nodeSt := make([]int, 0, len(to))
	edgeSt := make([]int, 0, len(to))
	// result of traversal
	resultEdges := make([]int, 0, cntege)
	// start Euler from node 1
	nodeSt = append(nodeSt, 1)
	edgeSt = append(edgeSt, -1)
	// iterative Hierholzer's algorithm
	for len(nodeSt) > 0 {
		v := nodeSt[len(nodeSt)-1]
		if iter[v] != -1 {
			e := iter[v]
			iter[v] = nxt[e]
			if vis[e] {
				continue
			}
			// mark edge visited
			vis[e] = true
			vis[e^1] = true
			// move to next node
			nodeSt = append(nodeSt, to[e])
			edgeSt = append(edgeSt, e)
		} else {
			// backtrack
			nodeSt = nodeSt[:len(nodeSt)-1]
			e := edgeSt[len(edgeSt)-1]
			edgeSt = edgeSt[:len(edgeSt)-1]
			if e != -1 {
				resultEdges = append(resultEdges, e)
			}
		}
	}
	// assign orientations alternately
	z := 0
	for _, e := range resultEdges {
		z++
		if z%2 == 1 {
			vise[e] = true
		} else {
			vise[e^1] = true
		}
	}
	// ensure even number of undirected edges
	if cntege%2 != 0 {
		// add one more undirected self-loop
		// add two directed edges without updating degrees
		// but count one undirected edge
		// first directed
		e1 := len(to)
		addDir(1, 1)
		// second directed
		e2 := len(to)
		addDir(1, 1)
		cntege++
		// mark orientation of the second directed edge
		vise[e2] = true
	}
	// output result
	fmt.Fprintln(writer, cntege)
	for u := 1; u <= n; u++ {
		for e := adj[u]; e != -1; e = nxt[e] {
			if vise[e] {
				fmt.Fprintln(writer, u, to[e])
			}
		}
	}
}
