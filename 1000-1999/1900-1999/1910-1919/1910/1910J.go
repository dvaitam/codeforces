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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	color := make([]int, n+1)
	redCount := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &color[i])
		if color[i] == 1 {
			redCount++
		}
	}
	blueCount := n - redCount

	adj := make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		var w int64
		fmt.Fscan(reader, &x, &y, &w)
		adj[x] = append(adj[x], Edge{y, w})
		adj[y] = append(adj[y], Edge{x, w})
	}

	if redCount != blueCount {
		fmt.Fprintln(writer, "Infinity")
		return
	}

	parent := make([]int, n+1)
	weight := make([]int64, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			weight[e.to] = e.w
			stack = append(stack, e.to)
		}
	}

	diff := make([]int64, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		if color[u] == 1 {
			diff[u] = 1
		} else {
			diff[u] = -1
		}
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			diff[u] += diff[e.to]
		}
	}

	var result int64
	for v := 2; v <= n; v++ {
		result += abs(diff[v]) * weight[v]
	}

	fmt.Fprintln(writer, result)
}
