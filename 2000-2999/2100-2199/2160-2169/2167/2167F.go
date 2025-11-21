package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		order := make([]int, 0, n)
		stack := []int{1}
		parent[1] = 0
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}

		sub := make([]int, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			sub[v] = 1
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				sub[v] += sub[to]
			}
		}

		ans := int64(n)
		for v := 2; v <= n; v++ {
			p := parent[v]
			if p == 0 {
				continue
			}
			sz := sub[v]
			if n-sz >= k {
				ans += int64(sz)
			}
			if sz >= k {
				ans += int64(n - sz)
			}
		}

		fmt.Fprintln(out, ans)
	}
}
