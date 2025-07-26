package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, adj [][]int) ([]int, int) {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, start)
	dist[start] = 0
	head := 0
	far := start
	for head < len(q) {
		v := q[head]
		head++
		if dist[v] > dist[far] {
			far = v
		}
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist, far
}

func approxDiameter(adj [][]int) int {
	_, v := bfs(0, adj)
	dist, _ := bfs(v, adj)
	d := 0
	for _, x := range dist {
		if x > d {
			d = x
		}
	}
	return d
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		if u >= 0 && v >= 0 && u < n && v < n {
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
	}

	ans := make([]int, 0, q+1)
	ans = append(ans, approxDiameter(adj))

	for i := 0; i < q; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			break
		}
		u--
		v--
		if u >= 0 && v >= 0 && u < n && v < n {
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		ans = append(ans, approxDiameter(adj))
	}

	for i, x := range ans {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, x)
	}
	fmt.Fprintln(writer)
}
