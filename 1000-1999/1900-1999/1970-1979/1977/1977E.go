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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		mat := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &mat[i])
		}

		// Build adjacency: edge from j to i if mat[i][j]=='1'
		adj := make([][]bool, n)
		for i := range adj {
			adj[i] = make([]bool, n)
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if mat[i][j] == '1' {
					adj[j][i] = true
				}
			}
		}

		// Transitive closure
		reach := make([][]bool, n)
		for i := range reach {
			reach[i] = make([]bool, n)
			copy(reach[i], adj[i])
			reach[i][i] = true
		}
		for k := 0; k < n; k++ {
			for i := 0; i < n; i++ {
				if reach[i][k] {
					for j := 0; j < n; j++ {
						if reach[k][j] {
							reach[i][j] = true
						}
					}
				}
			}
		}

		// Build bipartite graph for chain decomposition
		g := make([][]int, n)
		for j := 0; j < n; j++ {
			for i := 0; i < j; i++ {
				if reach[j][i] {
					g[j] = append(g[j], i)
				}
			}
		}

		// Hopcroft-Karp
		matchL := make([]int, n)
		matchR := make([]int, n)
		for i := range matchL {
			matchL[i] = -1
			matchR[i] = -1
		}
		INF := 1 << 30
		dist := make([]int, n)
		var bfs func() bool
		bfs = func() bool {
			q := make([]int, 0, n)
			for i := 0; i < n; i++ {
				if matchL[i] == -1 {
					dist[i] = 0
					q = append(q, i)
				} else {
					dist[i] = INF
				}
			}
			head := 0
			for head < len(q) {
				v := q[head]
				head++
				for _, u := range g[v] {
					w := matchR[u]
					if w != -1 && dist[w] == INF {
						dist[w] = dist[v] + 1
						q = append(q, w)
					}
				}
			}
			return true
		}
		var dfs func(int) bool
		dfs = func(v int) bool {
			for _, u := range g[v] {
				w := matchR[u]
				if w == -1 || (dist[w] == dist[v]+1 && dfs(w)) {
					matchL[v] = u
					matchR[u] = v
					return true
				}
			}
			dist[v] = INF
			return false
		}
		for {
			bfs()
			flow := 0
			for i := 0; i < n; i++ {
				if matchL[i] == -1 && dfs(i) {
					flow++
				}
			}
			if flow == 0 {
				break
			}
		}

		// Reconstruct chains and assign colors
		color := make([]int, n)
		for i := range color {
			color[i] = -1
		}
		chain := 0
		for v := 0; v < n; v++ {
			if matchR[v] == -1 {
				cur := v
				for cur != -1 {
					color[cur] = chain
					cur = matchL[cur]
				}
				chain++
			}
		}
		if chain == 1 { // all in one chain
			for i := 0; i < n; i++ {
				if i > 0 {
					out.WriteByte(' ')
				}
				out.WriteByte('0')
			}
			out.WriteByte('\n')
			continue
		}
		// chain should be <=2
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			if color[i] == 0 {
				out.WriteByte('0')
			} else {
				out.WriteByte('1')
			}
		}
		out.WriteByte('\n')
	}
}
