package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
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
		maxVal := make([]int, n+1)
		dp := make([]int, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			maxVal[u] = u
			heavy := 0
			for _, v := range g[u] {
				if v == parent[u] {
					continue
				}
				if maxVal[v] > maxVal[u] {
					maxVal[u] = maxVal[v]
				}
				if heavy == 0 || maxVal[v] > maxVal[heavy] {
					heavy = v
				}
			}
			if heavy == 0 {
				dp[u] = 1
				continue
			}
			res := dp[heavy]
			for _, v := range g[u] {
				if v == parent[u] || v == heavy {
					continue
				}
				res = res * (dp[v] + 1) % MOD
			}
			dp[u] = res % MOD
		}
		fmt.Fprintln(out, dp[1]%MOD)
	}
}
