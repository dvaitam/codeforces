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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	cost := make([][]int64, 3)
	for i := 0; i < 3; i++ {
		cost[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &cost[i][j])
		}
	}
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// tree must be a path: all degrees <=2
	start := -1
	for i := 0; i < n; i++ {
		if len(adj[i]) > 2 {
			fmt.Fprintln(writer, -1)
			return
		}
		if len(adj[i]) == 1 {
			start = i
		}
	}
	if start == -1 { // should not happen
		fmt.Fprintln(writer, -1)
		return
	}
	// build path order
	order := make([]int, n)
	prev := -1
	cur := start
	for i := 0; i < n; i++ {
		order[i] = cur
		next := -1
		for _, nb := range adj[cur] {
			if nb != prev {
				next = nb
				break
			}
		}
		prev = cur
		cur = next
	}
	perms := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}
	bestCost := int64(1<<63 - 1)
	bestColor := make([]int, n)
	for _, p := range perms {
		curCost := int64(0)
		colors := make([]int, n)
		for idx, node := range order {
			color := p[idx%3]
			curCost += cost[color][node]
			colors[node] = color + 1
		}
		if curCost < bestCost {
			bestCost = curCost
			copy(bestColor, colors)
		}
	}
	fmt.Fprintln(writer, bestCost)
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, bestColor[i])
	}
	writer.WriteByte('\n')
}
