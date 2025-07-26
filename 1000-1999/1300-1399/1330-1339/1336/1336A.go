package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	depth := make([]int, n+1)
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
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}

	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			size[u] += size[v]
		}
	}

	vals := make([]int, n)
	for i := 1; i <= n; i++ {
		vals[i-1] = depth[i] - (size[i] - 1)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })

	ans := 0
	for i := 0; i < k; i++ {
		ans += vals[i]
	}
	fmt.Fprintln(out, ans)
}
