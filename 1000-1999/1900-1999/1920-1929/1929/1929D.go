package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		parent := make([]int, n)
		order := make([]int, 0, n)
		stack := []int{0}
		parent[0] = -1
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
		dp := make([]int64, n)
		for i := n - 1; i >= 0; i-- {
			u := order[i]
			prod := int64(1)
			for _, v := range g[u] {
				if v == parent[u] {
					continue
				}
				prod = prod * (1 + dp[v]) % mod
			}
			dp[u] = prod
		}
		ans := int64(1) // empty set
		for _, val := range dp {
			ans += val
		}
		ans %= mod
		fmt.Fprintln(out, ans)
	}
}
