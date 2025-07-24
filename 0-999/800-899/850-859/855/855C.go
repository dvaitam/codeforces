package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

var (
	n   int
	m   int64
	adj [][]int
	k   int64
	x   int
	dp  [][][]int64
)

func dfs(u, p int) {
	// initialize base cases for node u
	dp[u][0][0] = (k - 1) % MOD
	if dp[u][0][0] < 0 {
		dp[u][0][0] += MOD
	}
	if x >= 1 {
		dp[u][1][1] = 1
	}
	dp[u][2][0] = (m - k) % MOD
	if dp[u][2][0] < 0 {
		dp[u][2][0] += MOD
	}

	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dfs(v, u)
		new0 := make([]int64, x+1)
		new1 := make([]int64, x+1)
		new2 := make([]int64, x+1)
		for i := 0; i <= x; i++ {
			if dp[u][0][i] != 0 {
				for j := 0; j+i <= x; j++ {
					// parent <k, child can be 0,1,2
					if dp[v][0][j] != 0 {
						new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][0][j]) % MOD
					}
					if dp[v][1][j] != 0 {
						new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][1][j]) % MOD
					}
					if dp[v][2][j] != 0 {
						new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][2][j]) % MOD
					}
				}
			}
			if dp[u][1][i] != 0 {
				for j := 0; j+i <= x; j++ {
					// parent k, child must be <k
					if dp[v][0][j] != 0 {
						new1[i+j] = (new1[i+j] + dp[u][1][i]*dp[v][0][j]) % MOD
					}
				}
			}
			if dp[u][2][i] != 0 {
				for j := 0; j+i <= x; j++ {
					// parent >k, child cannot be k
					if dp[v][0][j] != 0 {
						new2[i+j] = (new2[i+j] + dp[u][2][i]*dp[v][0][j]) % MOD
					}
					if dp[v][2][j] != 0 {
						new2[i+j] = (new2[i+j] + dp[u][2][i]*dp[v][2][j]) % MOD
					}
				}
			}
		}
		dp[u][0] = new0
		dp[u][1] = new1
		dp[u][2] = new2
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	fmt.Fscan(reader, &k, &x)

	dp = make([][][]int64, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([][]int64, 3)
		for s := 0; s < 3; s++ {
			dp[i][s] = make([]int64, x+1)
		}
	}

	dfs(1, 0)

	var ans int64
	for s := 0; s < 3; s++ {
		for c := 0; c <= x; c++ {
			ans = (ans + dp[1][s][c]) % MOD
		}
	}
	fmt.Fprintln(writer, ans)
}
