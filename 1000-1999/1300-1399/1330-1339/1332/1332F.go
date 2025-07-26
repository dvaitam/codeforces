package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}

	dp := make([][4]int64, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		prodA := int64(1)
		prodAB0 := int64(1)
		prodAplusB0B1 := int64(1)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			A := (dp[v][0] + dp[v][1]) % MOD
			B0 := dp[v][2]
			B1 := dp[v][3]
			prodA = prodA * A % MOD
			prodAB0 = prodAB0 * ((A + B0) % MOD) % MOD
			prodAplusB0B1 = prodAplusB0B1 * ((A + B0 + B1) % MOD) % MOD
		}
		dp[u][0] = prodAplusB0B1 % MOD
		dp[u][2] = dp[u][0]
		dp[u][3] = prodAB0 % MOD
		val := (prodAB0 - prodA) % MOD
		if val < 0 {
			val += MOD
		}
		dp[u][1] = val
	}
	ans := (dp[1][0] + dp[1][1] - 1) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Println(ans)
}
