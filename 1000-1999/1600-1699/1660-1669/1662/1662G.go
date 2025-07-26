package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct{ u, v int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 1; i < n; i++ {
		var p int
		fmt.Fscan(in, &p)
		p--
		adj[p] = append(adj[p], i)
		adj[i] = append(adj[i], p)
	}
	// BFS from node0 to find farthest
	bfs := func(s int) (int, []int) {
		dist := make([]int, n)
		for i := 0; i < n; i++ {
			dist[i] = -1
		}
		q := []int{s}
		dist[s] = 0
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		far := s
		for i := 0; i < n; i++ {
			if dist[i] > dist[far] {
				far = i
			}
		}
		return far, dist
	}
	s, _ := bfs(0)
	t, dist := bfs(s)
	_ = dist
	_ = t
	// BFS from s to assign labels
	order := make([]int, 0, n)
	q := []int{s}
	visited := make([]bool, n)
	visited[s] = true
	for head := 0; head < len(q); head++ {
		u := q[head]
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	label := make([]int, n)
	for idx, node := range order {
		label[node] = idx
	}
	// compute dp in reverse order
	// orient edges from smaller label to larger
	order2 := make([]int, n)
	copy(order2, order)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		order2[i], order2[j] = order2[j], order2[i]
	}
	dp := make([]int64, n)
	for _, u := range order2 {
		dp[u] = 1
		for _, v := range adj[u] {
			if label[u] < label[v] {
				dp[u] += dp[v]
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		ans += dp[i]
	}
	fmt.Println(ans)
}
