package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	g := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}

	visited := make([]bool, n+1)
	idx := make([]int, n+1)
	stack := []int{1}
	visited[1] = true
	order := make([]int, 0, 2*n)
	order = append(order, 1)
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		if idx[v] < len(g[v]) {
			u := g[v][idx[v]]
			idx[v]++
			if !visited[u] {
				visited[u] = true
				stack = append(stack, u)
				order = append(order, u)
			}
		} else {
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				order = append(order, stack[len(stack)-1])
			}
		}
	}

	t := (2*n + k - 1) / k
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	pos := 0
	for i := 0; i < k; i++ {
		if pos >= len(order) {
			fmt.Fprintln(out, "1 1")
		} else {
			remain := len(order) - pos
			take := t
			if take > remain {
				take = remain
			}
			fmt.Fprint(out, take)
			for j := 0; j < take; j++ {
				fmt.Fprint(out, " ", order[pos])
				pos++
			}
			fmt.Fprint(out, "\n")
		}
	}
}
