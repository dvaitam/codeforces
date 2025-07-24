package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

var g [][]int
var dp [][]int64

func dfs(u, p int) {
	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u)
	}
	// gather children
	children := []int{}
	for _, v := range g[u] {
		if v != p {
			children = append(children, v)
		}
	}
	m := len(children)

	if m == 0 {
		// leaf
		dp[u][0] = 0 // UpInternal needs one child
		dp[u][1] = 1 // UpEndpoint
		dp[u][2] = 1 // Isolated
		dp[u][3] = 0 // DownEndpoint
		dp[u][4] = 0 // Internal
		return
	}

	A := make([]int64, m)
	B := make([]int64, m)
	C := make([]int64, m)
	for i, v := range children {
		A[i] = (dp[v][0] + dp[v][1]) % MOD
		B[i] = (dp[v][2] + dp[v][3] + dp[v][4]) % MOD
		C[i] = dp[v][4] % MOD
	}

	// prefix and suffix for C
	preC := make([]int64, m+1)
	sufC := make([]int64, m+1)
	preC[0] = 1
	for i := 0; i < m; i++ {
		preC[i+1] = preC[i] * C[i] % MOD
	}
	sufC[m] = 1
	for i := m - 1; i >= 0; i-- {
		sufC[i] = sufC[i+1] * C[i] % MOD
	}

	prodC := preC[m]
	dp[u][1] = prodC
	dp[u][2] = prodC

	var down int64 = 0
	for i := 0; i < m; i++ {
		down = (down + A[i]*preC[i]%MOD*sufC[i+1]) % MOD
	}
	dp[u][3] = down

	// for UpInternal and Internal, use dynamic programming with B and A
	// dp0, dp1, dp2: ways keeping 0,1,2 edges so far
	var dp0, dp1, dp2 int64 = 1, 0, 0
	for i := 0; i < m; i++ {
		a := A[i]
		b := B[i]
		n0 := dp0 * b % MOD
		n1 := (dp1*b + dp0*a) % MOD
		n2 := (dp2*b + dp1*a) % MOD
		dp0, dp1, dp2 = n0, n1, n2
	}
	dp[u][0] = dp1
	dp[u][4] = dp2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dp = make([][]int64, n)
	for i := range dp {
		dp[i] = make([]int64, 5)
	}
	dfs(0, -1)
	ans := (dp[0][2] + dp[0][3] + dp[0][4]) % MOD
	fmt.Println(ans)
}
