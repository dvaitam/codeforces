package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k, d int
	if _, err := fmt.Fscan(in, &n, &k, &d); err != nil {
		return
	}
	police := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &police[i])
	}
	adj := make([][]Edge, n+1)
	for i := 1; i <= n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], Edge{v, i})
		adj[v] = append(adj[v], Edge{u, i})
	}

	visited := make([]bool, n+1)
	dist := make([]int, n+1)
	keep := make([]bool, n) // index 0..n-2 representing edges 1..n-1
	queue := make([]int, 0, n)
	for _, p := range police {
		if !visited[p] {
			visited[p] = true
			queue = append(queue, p)
		}
	}

	for head := 0; head < len(queue); head++ {
		u := queue[head]
		if dist[u] == d {
			continue
		}
		for _, e := range adj[u] {
			v := e.to
			if !visited[v] {
				visited[v] = true
				dist[v] = dist[u] + 1
				keep[e.idx-1] = true
				queue = append(queue, v)
			}
		}
	}

	result := make([]int, 0)
	for i := 1; i <= n-1; i++ {
		if !keep[i-1] {
			result = append(result, i)
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(result))
	for i, v := range result {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	if len(result) > 0 {
		out.WriteByte('\n')
	}
	out.Flush()
}
