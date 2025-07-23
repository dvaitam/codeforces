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

	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent := make([]int, n+1)
	depth := make([]int, n+1)
	stack := []int{1}
	parent[1] = 0
	depth[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}

	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(in, &u, &v)
		ans := 0
		dist := 0
		cur := v
		for cur != u {
			if val := a[cur] ^ dist; val > ans {
				ans = val
			}
			cur = parent[cur]
			dist++
		}
		if val := a[cur] ^ dist; val > ans {
			ans = val
		}
		fmt.Fprintln(out, ans)
	}
}
