package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 300005

var (
	g   [maxN][]int
	a   [maxN]int64
	dp  [maxN][20]int64
	rdr = bufio.NewReader(os.Stdin)
	wtr = bufio.NewWriter(os.Stdout)
)

func dfs(x, p int) {
	for i := 1; i < 20; i++ {
		dp[x][i] = int64(i) * a[x]
	}
	for _, u := range g[x] {
		if u == p {
			continue
		}
		dfs(u, x)
		for j := 1; j < 20; j++ {
			sum := int64(1<<63 - 1)
			for k := 1; k < 20; k++ {
				if j != k && dp[u][k] < sum {
					sum = dp[u][k]
				}
			}
			dp[x][j] += sum
		}
	}
}

func main() {
	defer wtr.Flush()
	var T int
	fmt.Fscan(rdr, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(rdr, &n)
		for i := 1; i <= n; i++ {
			fmt.Fscan(rdr, &a[i])
			g[i] = g[i][:0]
		}
		for i := 1; i < n; i++ {
			var x, y int
			fmt.Fscan(rdr, &x, &y)
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}
		dfs(1, 0)
		ans := dp[1][1]
		for i := 2; i < 20; i++ {
			if dp[1][i] < ans {
				ans = dp[1][i]
			}
		}
		fmt.Fprintln(wtr, ans)
	}
}
