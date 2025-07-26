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
	fmt.Fscan(in, &n, &m)
	if n < 5 {
		fmt.Fprintln(out, -1)
		return
	}
	limit := n
	if limit > 43 {
		limit = 43
	}
	adj := make([][]bool, limit)
	for i := 0; i < limit; i++ {
		adj[i] = make([]bool, limit)
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if u <= limit && v <= limit {
			u--
			v--
			adj[u][v] = true
			adj[v][u] = true
		}
	}
	ids := make([]int, limit)
	for i := 0; i < limit; i++ {
		ids[i] = i + 1
	}
	choose := make([]int, 5)
	var ans []int
	var dfs func(start, depth int)
	dfs = func(start, depth int) {
		if len(ans) > 0 {
			return
		}
		if depth == 5 {
			clique := true
			independent := true
			for i := 0; i < 5; i++ {
				for j := i + 1; j < 5; j++ {
					if adj[choose[i]-1][choose[j]-1] {
						independent = false
					} else {
						clique = false
					}
				}
			}
			if clique || independent {
				ans = append([]int(nil), choose...)
			}
			return
		}
		for i := start; i <= limit-(5-depth); i++ {
			choose[depth] = ids[i]
			dfs(i+1, depth+1)
			if len(ans) > 0 {
				return
			}
		}
	}
	dfs(0, 0)
	if len(ans) == 0 {
		fmt.Fprintln(out, -1)
	} else {
		for i := 0; i < 5; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
