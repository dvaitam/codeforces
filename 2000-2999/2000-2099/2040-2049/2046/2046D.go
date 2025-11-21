package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		messengers := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &messengers[i])
		}

		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
		}

		compID, compAdj := kosaraju(n, adj)
		sccCount := len(compAdj)
		if sccCount == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		nodesInComp := make([][]int, sccCount)
		hasMessenger := make([]bool, sccCount)
		for node := 0; node < n; node++ {
			cid := compID[node]
			nodesInComp[cid] = append(nodesInComp[cid], node)
			if messengers[node] > 0 {
				hasMessenger[cid] = true
			}
		}

		reachable := make([][]bool, sccCount)
		for i := 0; i < sccCount; i++ {
			reachable[i] = make([]bool, sccCount)
			dfsReach(i, i, compAdj, reachable[i])
		}

		best := inf
		for i := 0; i < sccCount; i++ {
			if !hasMessenger[i] {
				continue
			}
			all := true
			for j := 0; j < sccCount; j++ {
				if !reachable[i][j] {
					all = false
					break
				}
			}
			if all && len(nodesInComp[i]) < best {
				best = len(nodesInComp[i])
			}
		}

		if best == inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, best)
		}
	}
}

func dfsReach(start, v int, adj [][]int, seen []bool) {
	seen[v] = true
	for _, to := range adj[v] {
		if !seen[to] {
			dfsReach(start, to, adj, seen)
		}
	}
}

func kosaraju(n int, adj [][]int) ([]int, [][]int) {
	order := make([]int, 0, n)
	seen := make([]bool, n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		seen[v] = true
		for _, to := range adj[v] {
			if !seen[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < n; i++ {
		if !seen[i] {
			dfs1(i)
		}
	}

	radj := make([][]int, n)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			radj[v] = append(radj[v], u)
		}
	}

	compID := make([]int, n)
	for i := range compID {
		compID[i] = -1
	}
	var dfs2 func(int, int)
	dfs2 = func(v, cid int) {
		compID[v] = cid
		for _, to := range radj[v] {
			if compID[to] == -1 {
				dfs2(to, cid)
			}
		}
	}

	compCount := 0
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if compID[v] == -1 {
			dfs2(v, compCount)
			compCount++
		}
	}

	compAdj := make([][]int, compCount)
	exists := make(map[int]bool)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			cu := compID[u]
			cv := compID[v]
			if cu == cv {
				continue
			}
			key := cu*compCount + cv
			if !exists[key] {
				compAdj[cu] = append(compAdj[cu], cv)
				exists[key] = true
			}
		}
	}

	return compID, compAdj
}
