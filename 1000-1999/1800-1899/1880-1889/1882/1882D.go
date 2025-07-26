package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		parent := make([]int, n+1)
		size := make([]int, n+1)
		order := make([]int, 0, n)

		stack := []int{1}
		parent[1] = 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, v := range g[u] {
				if v != parent[u] {
					parent[v] = u
					stack = append(stack, v)
				}
			}
		}

		base := int64(0)
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			size[u] = 1
			for _, v := range g[u] {
				if v != parent[u] {
					size[u] += size[v]
					base += int64(size[v]) * int64(a[u]^a[v])
				}
			}
		}

		cost := make([]int64, n+1)
		cost[1] = base
		queue := []int{1}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range g[u] {
				if v == parent[u] {
					continue
				}
				w := a[u] ^ a[v]
				cost[v] = cost[u] + int64(n-2*size[v])*int64(w)
				queue = append(queue, v)
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, cost[i])
		}
		fmt.Fprintln(writer)
	}
}
