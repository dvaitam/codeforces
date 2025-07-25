package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		// Precompute parent and depth for LCA via DFS
		log := 0
		for (1 << log) <= n {
			log++
		}
		parent := make([][]int, log)
		for i := range parent {
			parent[i] = make([]int, n+1)
		}
		depth := make([]int, n+1)
		stack := []int{1}
		parent[0][1] = 0
		depth[1] = 0
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range adj[v] {
				if to == parent[0][v] {
					continue
				}
				parent[0][to] = v
				depth[to] = depth[v] + 1
				stack = append(stack, to)
			}
		}
		for k := 1; k < log; k++ {
			for i := 1; i <= n; i++ {
				parent[k][i] = parent[k-1][parent[k-1][i]]
			}
		}

		lca := func(a, b int) int {
			if depth[a] < depth[b] {
				a, b = b, a
			}
			diff := depth[a] - depth[b]
			for k := 0; diff > 0; k++ {
				if diff&1 == 1 {
					a = parent[k][a]
				}
				diff >>= 1
			}
			if a == b {
				return a
			}
			for k := log - 1; k >= 0; k-- {
				if parent[k][a] != parent[k][b] {
					a = parent[k][a]
					b = parent[k][b]
				}
			}
			return parent[0][a]
		}

		visited := make([]bool, n+1)
		used := make([]int, 0, n)

		for ; q > 0; q-- {
			aVals := make([]int, n+1)
			for i := 1; i <= n; i++ {
				fmt.Fscan(in, &aVals[i])
			}
			var u, v int
			fmt.Fscan(in, &u, &v)
			anc := lca(u, v)

			addVal := func(x int) {
				val := aVals[x]
				if !visited[val] {
					visited[val] = true
					used = append(used, val)
				}
			}

			for x := u; x != anc; x = parent[0][x] {
				addVal(x)
			}
			for x := v; x != anc; x = parent[0][x] {
				addVal(x)
			}
			addVal(anc)

			mex := 0
			for mex <= n && visited[mex] {
				mex++
			}
			fmt.Fprintln(out, mex)

			for _, val := range used {
				visited[val] = false
			}
			used = used[:0]
		}
	}
}
