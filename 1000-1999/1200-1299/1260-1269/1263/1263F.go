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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	// Main grid
	var a int
	fmt.Fscan(reader, &a)
	parA := make([]int, a+1) // 1-indexed; parA[1]=0
	for i := 2; i <= a; i++ {
		fmt.Fscan(reader, &parA[i])
	}
	x := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &x[i])
	}

	// Reserve grid
	var b int
	fmt.Fscan(reader, &b)
	parB := make([]int, b+1)
	for i := 2; i <= b; i++ {
		fmt.Fscan(reader, &parB[i])
	}
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &y[i])
	}

	// Build trees adjacency
	chA := make([][]int, a+1)
	for i := 2; i <= a; i++ {
		p := parA[i]
		chA[p] = append(chA[p], i)
	}
	chB := make([][]int, b+1)
	for i := 2; i <= b; i++ {
		p := parB[i]
		chB[p] = append(chB[p], i)
	}

	leafIdxA := make([]int, a+1)
	leafIdxB := make([]int, b+1)
	for i := 1; i <= n; i++ {
		leafIdxA[x[i]] = i
		leafIdxB[y[i]] = i
	}

	// compute intervals via DFS
	lA := make([]int, a+1)
	rA := make([]int, a+1)
	var dfsA func(int)
	dfsA = func(u int) {
		if len(chA[u]) == 0 {
			idx := leafIdxA[u]
			lA[u], rA[u] = idx, idx
			return
		}
		l, r := n+1, 0
		for _, v := range chA[u] {
			dfsA(v)
			if lA[v] < l {
				l = lA[v]
			}
			if rA[v] > r {
				r = rA[v]
			}
		}
		lA[u], rA[u] = l, r
	}
	dfsA(1)

	lB := make([]int, b+1)
	rB := make([]int, b+1)
	var dfsB func(int)
	dfsB = func(u int) {
		if len(chB[u]) == 0 {
			idx := leafIdxB[u]
			lB[u], rB[u] = idx, idx
			return
		}
		l, r := n+1, 0
		for _, v := range chB[u] {
			dfsB(v)
			if lB[v] < l {
				l = lB[v]
			}
			if rB[v] > r {
				r = rB[v]
			}
		}
		lB[u], rB[u] = l, r
	}
	dfsB(1)

	// build list of edges (excluding roots)
	type edge struct {
		l, r int
	}
	edgesA := make([]edge, 0, a-1)
	for i := 2; i <= a; i++ {
		edgesA = append(edgesA, edge{lA[i], rA[i]})
	}
	edgesB := make([]edge, 0, b-1)
	for i := 2; i <= b; i++ {
		edgesB = append(edgesB, edge{lB[i], rB[i]})
	}

	n1 := len(edgesA)
	n2 := len(edgesB)
	adj := make([][]int, n1)
	for i := 0; i < n1; i++ {
		e1 := edgesA[i]
		for j := 0; j < n2; j++ {
			e2 := edgesB[j]
			if e1.l <= e2.r && e2.l <= e1.r {
				adj[i] = append(adj[i], j)
			}
		}
	}

	// Hopcroft-Karp
	pairU := make([]int, n1)
	for i := range pairU {
		pairU[i] = -1
	}
	pairV := make([]int, n2)
	for i := range pairV {
		pairV[i] = -1
	}
	dist := make([]int, n1)
	INF := int(1e9)

	bfs := func() bool {
		q := []int{}
		for i := 0; i < n1; i++ {
			if pairU[i] == -1 {
				dist[i] = 0
				q = append(q, i)
			} else {
				dist[i] = INF
			}
		}
		found := false
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range adj[u] {
				pu := pairV[v]
				if pu != -1 && dist[pu] == INF {
					dist[pu] = dist[u] + 1
					q = append(q, pu)
				}
				if pu == -1 {
					found = true
				}
			}
		}
		return found
	}

	var dfs func(int) bool
	dfs = func(u int) bool {
		for _, v := range adj[u] {
			pu := pairV[v]
			if pu == -1 || (dist[pu] == dist[u]+1 && dfs(pu)) {
				pairU[u] = v
				pairV[v] = u
				return true
			}
		}
		dist[u] = INF
		return false
	}

	matching := 0
	for bfs() {
		for i := 0; i < n1; i++ {
			if pairU[i] == -1 {
				if dfs(i) {
					matching++
				}
			}
		}
	}

	totalEdges := (a - 1) + (b - 1)
	result := totalEdges - matching
	fmt.Fprintln(writer, result)
}
