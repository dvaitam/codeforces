package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(adj [][]int, start, banned int) []bool {
	vis := make([]bool, len(adj))
	queue := []int{start}
	vis[start] = true
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range adj[v] {
			if to == banned || vis[to] {
				continue
			}
			vis[to] = true
			queue = append(queue, to)
		}
	}
	return vis
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, a, b int
		fmt.Fscan(in, &n, &m, &a, &b)
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		visA := bfs(adj, a, b)
		visB := bfs(adj, b, a)
		cntA, cntB := 0, 0
		for i := 1; i <= n; i++ {
			if visA[i] {
				cntA++
			}
			if visB[i] {
				cntB++
			}
		}
		sizeA := n - cntA - 1
		sizeB := n - cntB - 1
		fmt.Fprintln(out, int64(sizeA)*int64(sizeB))
	}
}
