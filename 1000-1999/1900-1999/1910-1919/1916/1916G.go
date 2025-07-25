package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

var g [][]Edge
var ans int64

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func dfs(u, p int) map[int64]int {
	dp := make(map[int64]int)
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		child := dfs(e.to, u)
		for g1, l1 := range dp {
			for g2, l2 := range child {
				gp := gcd(g1, gcd(g2, e.w))
				l := l1 + l2 + 1
				if val := int64(l) * gp; val > ans {
					ans = val
				}
			}
		}
		for g2, l2 := range child {
			gNew := gcd(g2, e.w)
			if l2+1 > dp[gNew] {
				dp[gNew] = l2 + 1
			}
			if val := int64(l2+1) * gNew; val > ans {
				ans = val
			}
		}
		if dp[e.w] < 1 {
			dp[e.w] = 1
		}
		if e.w > ans {
			ans = e.w
		}
	}
	return dp
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
		g = make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			var w int64
			fmt.Fscan(in, &u, &v, &w)
			g[u] = append(g[u], Edge{v, w})
			g[v] = append(g[v], Edge{u, w})
		}
		ans = 0
		dfs(1, 0)
		fmt.Fprintln(out, ans)
	}
}
