package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

var (
	n, k int
	adj  [][]int
)

func dfs(v, p int) []int {
	dp := make([]int, k+1)
	dp[0] = 1
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		child := dfs(to, v)
		sumChild := 0
		for _, x := range child {
			sumChild += x
			if sumChild >= mod {
				sumChild -= mod
			}
		}
		newDp := make([]int, k+1)
		// cut edge v-to
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			val := dp[i] * sumChild % mod
			newDp[i] = (newDp[i] + val) % mod
		}
		// keep edge v-to
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := 0; j <= k; j++ {
				if child[j] == 0 {
					continue
				}
				if i+j+1 > k {
					continue
				}
				nd := i
				if j+1 > nd {
					nd = j + 1
				}
				val := (dp[i] * child[j]) % mod
				newDp[nd] += val
				if newDp[nd] >= mod {
					newDp[nd] -= mod
				}
			}
		}
		dp = newDp
	}
	return dp
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	adj = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	dp := dfs(0, -1)
	ans := 0
	for _, x := range dp {
		ans += x
		if ans >= mod {
			ans -= mod
		}
	}
	fmt.Fprintln(writer, ans%mod)
}
