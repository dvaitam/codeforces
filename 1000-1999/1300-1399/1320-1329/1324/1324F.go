package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x == 1 {
			val[i] = 1
		} else {
			val[i] = -1
		}
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := make([]int, 0, n)
	stack = append(stack, 1)
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stack = append(stack, u)
		}
	}

	dp := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		dp[v] = val[v]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			if dp[u] > 0 {
				dp[v] += dp[u]
			}
		}
	}

	up := make([]int, n+1)
	res := make([]int, n+1)
	stack = stack[:0]
	stack = append(stack, 1)
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res[v] = dp[v] + up[v]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			tmp := up[v] + dp[v] - max(0, dp[u])
			if tmp < 0 {
				tmp = 0
			}
			up[u] = tmp
			stack = append(stack, u)
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
	writer.Flush()
}
