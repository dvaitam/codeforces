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
		var n, m int
		var c int64
		fmt.Fscan(in, &n, &m, &c)
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		order := make([]int, 0, n)
		stack := []int{0}
		parent[0] = -2
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}
		size := make([]int, n)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			size[v] = 1
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				size[v] += size[to]
			}
		}

		minCost := int64(^uint64(0) >> 1)
		for v := 1; v < n; v++ {
			x := size[v]
			y := n - x
			cost := int64(x*x + y*y)
			if cost < minCost {
				minCost = cost
			}
		}
		fmt.Fprintln(out, minCost)
	}
}
