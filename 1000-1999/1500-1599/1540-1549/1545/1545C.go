package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		rows := make([][]int, 2*n)
		for i := 0; i < 2*n; i++ {
			rows[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &rows[i][j])
			}
		}
		rowDone := make([]bool, 2*n)
		must := make([]bool, 2*n)
		done := make([][]bool, n)
		for i := 0; i < n; i++ {
			done[i] = make([]bool, n+1)
		}
		cnt := make([]int, n+1)
		occ := make([][]int, n+1)
		adj := make([][]int, 2*n)
		vis := make([]int, 2*n)
		ans := make([]int, 0, n)
		ways := 1

		for it := 0; it < n; it++ {
			col, val := -1, -1
			// find unique
			for c := 0; c < n && col == -1; c++ {
				for i := 1; i <= n; i++ {
					cnt[i] = 0
				}
				for r := 0; r < 2*n; r++ {
					if !rowDone[r] {
						cnt[rows[r][c]]++
					}
				}
				for i := 1; i <= n; i++ {
					if done[c][i] {
						continue
					}
					if cnt[i] == 1 {
						col = c
						val = i
						break
					}
				}
			}
			if col != -1 {
				// take the row with this unique value
				var mr int
				for r := 0; r < 2*n; r++ {
					if !rowDone[r] && rows[r][col] == val {
						mr = r
						break
					}
				}
				ans = append(ans, mr)
				must[mr] = true
				// mark done values
				for c := 0; c < n; c++ {
					done[c][rows[mr][c]] = true
				}
				// mark conflicting rows
				for r := 0; r < 2*n; r++ {
					if rowDone[r] {
						continue
					}
					for c := 0; c < n; c++ {
						if rows[r][c] == rows[mr][c] {
							rowDone[r] = true
							break
						}
					}
				}
				continue
			}
			// no unique: build graph
			for i := range adj {
				adj[i] = adj[i][:0]
			}
			for c := 0; c < n; c++ {
				for i := 1; i <= n; i++ {
					cnt[i] = 0
					occ[i] = occ[i][:0]
				}
				for r := 0; r < 2*n; r++ {
					if !rowDone[r] {
						v := rows[r][c]
						cnt[v]++
						occ[v] = append(occ[v], r)
					}
				}
				for i := 1; i <= n; i++ {
					if done[c][i] {
						continue
					}
					// occ[i] has exactly two
					u1, u2 := occ[i][0], occ[i][1]
					adj[u1] = append(adj[u1], u2)
					adj[u2] = append(adj[u2], u1)
				}
			}
			for i := 0; i < 2*n; i++ {
				vis[i] = -1
			}
			// propagate fixed
			var dfs func(u, colr int)
			dfs = func(u, colr int) {
				if vis[u] != -1 {
					return
				}
				vis[u] = colr
				if rowDone[u] {
					// must match
				}
				for _, v := range adj[u] {
					dfs(v, 1-colr)
				}
			}
			for r := 0; r < 2*n; r++ {
				if rowDone[r] && vis[r] == -1 {
					colr := 0
					if must[r] {
						colr = 1
					}
					dfs(r, colr)
				}
			}
			for r := 0; r < 2*n; r++ {
				if vis[r] == -1 {
					dfs(r, 0)
					ways = ways * 2 % MOD
				}
			}
			ans = ans[:0]
			for r := 0; r < 2*n; r++ {
				if vis[r] == 1 {
					ans = append(ans, r)
				}
			}
			break
		}
		fmt.Fprintln(writer, ways)
		for i, r := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", r+1)
		}
		writer.WriteByte('\n')
	}
}
