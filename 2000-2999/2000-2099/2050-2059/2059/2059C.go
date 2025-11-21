package main

import (
	"bufio"
	"fmt"
	"os"
)

func dfs(v int, adj [][]int, matchQueue []int, visited []bool) bool {
	for _, q := range adj[v] {
		if visited[q] {
			continue
		}
		visited[q] = true
		if matchQueue[q] == -1 || dfs(matchQueue[q], adj, matchQueue, visited) {
			matchQueue[q] = v
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}

		suffix := make([][]int64, n)
		for i := 0; i < n; i++ {
			suffix[i] = make([]int64, n+1)
			for j := n - 1; j >= 0; j-- {
				suffix[i][j] = suffix[i][j+1] + a[i][j]
			}
		}

		adj := make([][]int, n)
		for v := 0; v < n; v++ {
			tIdx := n - v
			if tIdx < 0 || tIdx > n {
				continue
			}
			for i := 0; i < n; i++ {
				if suffix[i][tIdx] == int64(v) {
					adj[v] = append(adj[v], i)
				}
			}
		}

		matchQueue := make([]int, n)
		for i := range matchQueue {
			matchQueue[i] = -1
		}

		ans := 0
		for v := 0; v < n; v++ {
			if len(adj[v]) == 0 {
				break
			}
			visited := make([]bool, n)
			if dfs(v, adj, matchQueue, visited) {
				ans++
			} else {
				break
			}
		}

		fmt.Fprintln(out, ans)
	}
}
