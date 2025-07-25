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

	var n, t int
	if _, err := fmt.Fscan(reader, &n, &t); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var start int
	fmt.Fscan(reader, &start)

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{start}
	parent[start] = 0
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

	grundy := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		set := map[int]bool{}
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			set[grundy[v]] = true
		}
		g := 0
		for {
			if !set[g] {
				break
			}
			g++
		}
		grundy[u] = g
	}

	if grundy[start] != 0 {
		fmt.Fprintln(writer, "Ron")
	} else {
		fmt.Fprintln(writer, "Hermione")
	}
}
