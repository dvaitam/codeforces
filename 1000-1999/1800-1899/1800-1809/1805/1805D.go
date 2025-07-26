package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func bfs(start int, adj [][]int) []int {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// find diameter endpoints
	d1 := bfs(1, adj)
	u := 1
	for i := 1; i <= n; i++ {
		if d1[i] > d1[u] {
			u = i
		}
	}
	du := bfs(u, adj)
	v := u
	for i := 1; i <= n; i++ {
		if du[i] > du[v] {
			v = i
		}
	}
	dv := bfs(v, adj)

	f := make([]int, n)
	for i := 1; i <= n; i++ {
		if du[i] > dv[i] {
			f[i-1] = du[i]
		} else {
			f[i-1] = dv[i]
		}
	}

	sort.Ints(f)
	ans := make([]int, n)
	idx := 0
	for k := 1; k <= n; k++ {
		for idx < n && f[idx] < k {
			idx++
		}
		if idx == n {
			ans[k-1] = n
		} else {
			ans[k-1] = idx + 1
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
