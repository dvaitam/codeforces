package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	// choose any leaf as root
	root := 0
	for i := 0; i < n; i++ {
		if len(adj[i]) == 1 {
			root = i
			break
		}
	}

	depth := make([]int, n)
	visited := make([]bool, n)
	var dfs func(int, int)
	dfs = func(u int, d int) {
		visited[u] = true
		depth[u] = d
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v, d+1)
			}
		}
	}
	dfs(root, 0)

	parity0 := false
	parity1 := false
	for i := 0; i < n; i++ {
		if len(adj[i]) == 1 { // leaf
			if depth[i]%2 == 0 {
				parity0 = true
			} else {
				parity1 = true
			}
		}
	}
	minVal := 1
	if parity0 && parity1 {
		minVal = 3
	}

	maxVal := n - 1
	for u := 0; u < n; u++ {
		leafChildren := 0
		for _, v := range adj[u] {
			if len(adj[v]) == 1 {
				leafChildren++
			}
		}
		if leafChildren > 1 {
			maxVal -= leafChildren - 1
		}
	}

	fmt.Printf("%d %d\n", minVal, maxVal)
}
