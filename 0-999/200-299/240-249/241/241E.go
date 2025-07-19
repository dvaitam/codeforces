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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	u := make([]int, m)
	v := make([]int, m)
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &u[i], &v[i])
		u[i]--
		v[i]--
		adj[u[i]] = append(adj[u[i]], v[i])
	}
	vis := make([]bool, n)
	reach := make([]bool, n)
	var dfs func(int)
	dfs = func(x int) {
		vis[x] = true
		if x == n-1 {
			reach[x] = true
		}
		for _, w := range adj[x] {
			if !vis[w] {
				dfs(w)
			}
			if reach[w] {
				reach[x] = true
			}
		}
	}
	dfs(0)
	s := make([]int, n)
	updated := true
	for l := 1; l <= n && updated; l++ {
		updated = false
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if reach[ui] && reach[vi] {
				diff := s[ui] - s[vi]
				if diff < 1 {
					s[vi] = s[ui] - 1
					updated = true
				} else if diff > 2 {
					s[ui] = s[vi] + 2
					updated = true
				}
			}
		}
	}
	if updated {
		fmt.Fprintln(out, "No")
	} else {
		fmt.Fprintln(out, "Yes")
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if reach[ui] && reach[vi] {
				fmt.Fprintln(out, s[ui]-s[vi])
			} else {
				fmt.Fprintln(out, 1)
			}
		}
	}
}
