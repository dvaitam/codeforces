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

	var n, m int
	fmt.Fscan(reader, &n, &m)

	s := make([]int, n+1)
	a := make([][]int, n+1)
	cnt := 0
	root := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &s[i])
		if s[i] > cnt {
			cnt = s[i]
			root = i
		}
		if s[i] > 0 {
			a[i] = make([]int, s[i])
			for j := 0; j < s[i]; j++ {
				fmt.Fscan(reader, &a[i][j])
			}
		}
	}
	// build tree
	g := make([][]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	used := make([]int, m+2)
	col := make([]int, m+2)
	parent := make([]int, n+1)

	// iterative DFS from root
	stack := make([]int, 0, n)
	stack = append(stack, root)
	parent[root] = 0
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// mark used colors for queries at x
		for _, y := range a[x] {
			if col[y] != 0 {
				used[col[y]] = x
			}
		}
		// assign colors
		j := 1
		for _, y := range a[x] {
			if col[y] != 0 {
				continue
			}
			for used[j] == x {
				j++
			}
			col[y] = j
			j++
		}
		// traverse children
		for _, y := range g[x] {
			if y != parent[x] {
				parent[y] = x
				stack = append(stack, y)
			}
		}
	}

	if cnt == 0 {
		cnt = 1
	}
	fmt.Fprintln(writer, cnt)
	for i := 1; i <= m; i++ {
		if col[i] == 0 {
			col[i] = 1
		}
		fmt.Fprintf(writer, "%d ", col[i])
	}
	writer.WriteByte('\n')
}
