package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func bfs(start int, adj [][]int) []int {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, n)
	head, tail := 0, 0
	dist[start] = 0
	queue[tail] = start
	tail++
	for head < tail {
		v := queue[head]
		head++
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue[tail] = to
				tail++
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(i, adj)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			pathLen := dist[i][j]
			onPath := make([]bool, n)
			ok := true
			for d := 0; d <= pathLen; d++ {
				cnt := 0
				var who int
				for v := 0; v < n; v++ {
					if dist[i][v] == d && dist[j][v] == pathLen-d {
						cnt++
						who = v
					}
				}
				if cnt != 1 {
					ok = false
					break
				}
				onPath[who] = true
			}
			if !ok {
				fmt.Fprint(writer, 0)
				if j+1 == n {
					fmt.Fprintln(writer)
				} else {
					fmt.Fprint(writer, " ")
				}
				continue
			}
			ans := 1
			for v := 0; v < n; v++ {
				if onPath[v] {
					continue
				}
				cnt := 0
				for _, to := range adj[v] {
					if dist[i][to] == dist[i][v]-1 && dist[j][to] == dist[j][v]-1 {
						cnt++
					}
				}
				ans = ans * cnt % mod
				if cnt == 0 {
					break
				}
			}
			fmt.Fprint(writer, ans)
			if j+1 == n {
				fmt.Fprintln(writer)
			} else {
				fmt.Fprint(writer, " ")
			}
		}
	}
}
