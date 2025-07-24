package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m  int
	edges [][2]int
	best  int
)

func dfs(v int, color []int) {
	if v == n {
		// evaluate this coloring
		used := [7][7]bool{}
		count := 0
		for _, e := range edges {
			a := color[e[0]]
			b := color[e[1]]
			if a > b {
				a, b = b, a
			}
			if !used[a][b] {
				used[a][b] = true
				count++
			}
		}
		if count > best {
			best = count
		}
		return
	}
	for c := 1; c <= 6; c++ {
		color[v] = c
		dfs(v+1, color)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges = make([][2]int, m)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		edges[i] = [2]int{a - 1, b - 1}
	}
	color := make([]int, n)
	dfs(0, color)
	fmt.Fprintln(writer, best)
}
