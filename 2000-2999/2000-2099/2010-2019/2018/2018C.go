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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		depth := make([]int, n+1)
		order := make([]int, 0, n)
		stack := []int{1}
		parent[1] = 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				depth[v] = depth[u] + 1
				stack = append(stack, v)
			}
		}

		maxDepth := 0
		for i := 1; i <= n; i++ {
			if depth[i] > maxDepth {
				maxDepth = depth[i]
			}
		}

		mxDepth := make([]int, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			mx := depth[u]
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				if mxDepth[v] > mx {
					mx = mxDepth[v]
				}
			}
			mxDepth[u] = mx
		}

		if maxDepth == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		diff := make([]int, maxDepth+3)
		for u := 1; u <= n; u++ {
			left := depth[u]
			if left < 1 {
				left = 1
			}
			right := mxDepth[u]
			if right < left {
				continue
			}
			diff[left]++
			if right+1 < len(diff) {
				diff[right+1]--
			}
		}

		best := 0
		cur := 0
		for d := 1; d <= maxDepth; d++ {
			cur += diff[d]
			if cur > best {
				best = cur
			}
		}

		if best > n {
			best = n
		}
		fmt.Fprintln(out, n-best)
	}
}
