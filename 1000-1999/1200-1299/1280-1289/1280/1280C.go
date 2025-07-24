package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k int
		if _, err := fmt.Fscan(in, &k); err != nil {
			return
		}
		n := 2 * k
		g := make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var a, b int
			var w int64
			fmt.Fscan(in, &a, &b, &w)
			g[a] = append(g[a], Edge{b, w})
			g[b] = append(g[b], Edge{a, w})
		}

		parent := make([]int, n+1)
		weight := make([]int64, n+1) // edge weight from node to parent
		order := make([]int, 0, n)
		stack := []int{1}
		parent[1] = -1
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, e := range g[u] {
				if e.to == parent[u] {
					continue
				}
				parent[e.to] = u
				weight[e.to] = e.w
				stack = append(stack, e.to)
			}
		}

		size := make([]int, n+1)
		var G, B int64
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			size[u]++
			if parent[u] != -1 {
				s := size[u]
				if s%2 == 1 {
					G += weight[u]
				}
				if s < n-s {
					B += int64(s) * weight[u]
				} else {
					B += int64(n-s) * weight[u]
				}
				size[parent[u]] += s
			}
		}
		fmt.Fprintln(out, G, B)
	}
}
