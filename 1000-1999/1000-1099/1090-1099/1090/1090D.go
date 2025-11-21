package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([]map[int]struct{}, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if adj[a] == nil {
			adj[a] = make(map[int]struct{})
		}
		if adj[b] == nil {
			adj[b] = make(map[int]struct{})
		}
		adj[a][b] = struct{}{}
		adj[b][a] = struct{}{}
	}

	u := -1
	for i := 1; i <= n; i++ {
		deg := 0
		if adj[i] != nil {
			deg = len(adj[i])
		}
		if deg < n-1 {
			u = i
			break
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if u == -1 {
		fmt.Fprintln(out, "NO")
		return
	}

	marked := make([]bool, n+1)
	marked[u] = true
	if adj[u] != nil {
		for v := range adj[u] {
			marked[v] = true
		}
	}

	v := -1
	for i := 1; i <= n; i++ {
		if !marked[i] {
			v = i
			break
		}
	}

	if v == -1 {
		fmt.Fprintln(out, "NO")
		return
	}

	order := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if i != u && i != v {
			order = append(order, i)
		}
	}
	order = append(order, u, v)

	first := make([]int, n)
	second := make([]int, n)
	for idx, node := range order {
		first[node-1] = idx + 1
		if idx < n-1 {
			second[node-1] = idx + 1
		} else {
			second[node-1] = n - 1
		}
	}

	fmt.Fprintln(out, "YES")
	printArray(out, first)
	printArray(out, second)
}

func printArray(out *bufio.Writer, arr []int) {
	for i, val := range arr {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
	}
	fmt.Fprintln(out)
}
