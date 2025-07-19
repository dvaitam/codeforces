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
	G := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		G[x] = append(G[x], y)
		G[y] = append(G[y], x)
	}
	dist1, memo1 := bfs(n, G, 1)
	dist2, memo2 := bfs(n, G, n)

	total := memo1[n]
	best := 0.0
	shortest := dist1[n]
	for i := 1; i <= n; i++ {
		if dist1[i]+dist2[i] != shortest {
			continue
		}
		safe := memo1[i] * memo2[i]
		cur := float64(safe) / float64(total)
		if i != 1 && i != n {
			cur *= 2.0
		}
		if cur > best {
			best = cur
		}
	}
	fmt.Fprintf(writer, "%.9f\n", best)
}

func bfs(n int, G [][]int, start int) ([]int, []int64) {
	inf := int(1e9)
	dist := make([]int, n+1)
	memo := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	queue := make([]int, 0, n)
	head := 0
	dist[start] = 0
	memo[start] = 1
	queue = append(queue, start)
	for head < len(queue) {
		u := queue[head]
		head++
		for _, v := range G[u] {
			if dist[v] == inf {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
				memo[v] += memo[u]
			} else if dist[v] == dist[u]+1 {
				memo[v] += memo[u]
			}
		}
	}
	return dist, memo
}
