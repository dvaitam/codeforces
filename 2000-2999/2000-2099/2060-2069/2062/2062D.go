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

		l := make([]int64, n)
		r := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
		}

		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		if n == 1 {
			fmt.Fprintln(out, l[0])
			continue
		}

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		order := make([]int, 0, n)
		stack := []int{0}
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				stack = append(stack, v)
			}
		}

		T := make([]int64, n)
		C := make([]int64, n)

		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			need := l[u]
			sumC := int64(0)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				if T[v] > need {
					need = T[v]
				}
				sumC += C[v]
			}
			if need > r[u] {
				need = r[u]
			}
			extra := int64(0)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				if T[v] > need {
					extra += T[v] - need
				}
			}
			T[u] = need
			C[u] = sumC + extra
		}

		fmt.Fprintln(out, T[0]+C[0])
	}
}
