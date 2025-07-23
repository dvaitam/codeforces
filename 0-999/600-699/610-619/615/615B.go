package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	deg := make([]int, n+1)
	// neighbors with smaller index for DP
	prev := make([][]int, n+1)

	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		deg[u]++
		deg[v]++
		if u < v {
			prev[v] = append(prev[v], u)
		} else {
			prev[u] = append(prev[u], v)
		}
	}

	dp := make([]int, n+1)
	ans := 0
	for i := 1; i <= n; i++ {
		dp[i] = 1
		for _, p := range prev[i] {
			if dp[p]+1 > dp[i] {
				dp[i] = dp[p] + 1
			}
		}
		beauty := dp[i] * deg[i]
		if beauty > ans {
			ans = beauty
		}
	}

	fmt.Fprintln(writer, ans)
}
