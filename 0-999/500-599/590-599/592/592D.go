package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m    int
	adj     [][]int
	special []bool
	cnt     int
)

func dfs(u, p int) bool {
	need := special[u]
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		if dfs(v, u) {
			cnt++
			need = true
		}
	}
	return need
}

func bfs(start int) (int, []int) {
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	far := start
	maxd := 0
	for head := 0; head < len(q); head++ {
		u := q[head]
		if special[u] {
			if dist[u] > maxd || (dist[u] == maxd && u < far) {
				far = u
				maxd = dist[u]
			}
		}
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return far, dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n, &m)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	special = make([]bool, n+1)
	specials := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &specials[i])
		special[specials[i]] = true
	}

	cnt = 0
	dfs(specials[0], 0)

	far1, _ := bfs(specials[0])
	far2, dist2 := bfs(far1)
	d := dist2[far2]
	startCity := far1
	if far2 < startCity {
		startCity = far2
	}
	time := 2*cnt - d
	fmt.Fprintln(out, startCity)
	fmt.Fprintln(out, time)
}
