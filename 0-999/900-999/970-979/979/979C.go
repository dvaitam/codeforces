package main

import (
	"bufio"
	"fmt"
	"os"
)

func dfsInfo(root int, adj [][]int) ([]int, []int) {
	n := len(adj) - 1
	parent := make([]int, n+1)
	size := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{root}
	parent[root] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			size[v] += size[to]
		}
	}
	return parent, size
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, x, y int
	if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	parentX, sizeX := dfsInfo(x, adj)
	parentY, sizeY := dfsInfo(y, adj)

	xp := y
	for parentX[xp] != x {
		xp = parentX[xp]
	}
	yq := x
	for parentY[yq] != y {
		yq = parentY[yq]
	}

	sizeA := n - sizeX[xp]
	sizeB := n - sizeY[yq]
	totalPairs := int64(n) * int64(n-1)
	invalid := int64(sizeA) * int64(sizeB)
	ans := totalPairs - invalid
	fmt.Fprintln(writer, ans)
}
