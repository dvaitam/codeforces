package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(in, &n, &a, &b)
		g := make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			g[u] = append(g[u], Edge{v, w})
			g[v] = append(g[v], Edge{u, w})
		}

		// xor from a to all nodes
		xorA := make([]int, n+1)
		visited := make([]bool, n+1)
		stack := []int{a}
		visited[a] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, e := range g[v] {
				if !visited[e.to] {
					visited[e.to] = true
					xorA[e.to] = xorA[v] ^ e.w
					stack = append(stack, e.to)
				}
			}
		}

		// if direct path from a to b has xor 0
		if xorA[b] == 0 {
			fmt.Fprintln(out, "YES")
			continue
		}

		// xor from b to all nodes
		xorB := make([]int, n+1)
		visited = make([]bool, n+1)
		stack = []int{b}
		visited[b] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, e := range g[v] {
				if !visited[e.to] {
					visited[e.to] = true
					xorB[e.to] = xorB[v] ^ e.w
					stack = append(stack, e.to)
				}
			}
		}

		// collect xor values reachable from a without entering b
		reachable := make(map[int]struct{})
		visited = make([]bool, n+1)
		stack = []int{a}
		visited[a] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			reachable[xorA[v]] = struct{}{}
			for _, e := range g[v] {
				if e.to == b || visited[e.to] {
					continue
				}
				visited[e.to] = true
				stack = append(stack, e.to)
			}
		}

		// collect xorB for all nodes except b
		valsB := make(map[int]struct{})
		for i := 1; i <= n; i++ {
			if i == b {
				continue
			}
			valsB[xorB[i]] = struct{}{}
		}

		ans := "NO"
		for v := range reachable {
			if _, ok := valsB[v]; ok {
				ans = "YES"
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
