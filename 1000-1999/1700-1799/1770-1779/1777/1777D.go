package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
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

		height := make([]int, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			h := 0
			for _, v := range g[u] {
				if v == parent[u] {
					continue
				}
				if height[v]+1 > h {
					h = height[v] + 1
				}
			}
			height[u] = h
		}

		var total int64
		for i := 1; i <= n; i++ {
			total += int64(height[i] + 1)
		}
		pow := modPow(2, int64(n-1))
		ans := total % MOD * pow % MOD
		fmt.Fprintln(out, ans)
	}
}
