package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	cost := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cost[i])
	}
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	visited := make([]bool, n)
	q := make([]int, 0)
	var ans int64
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		visited[i] = true
		q = append(q[:0], i)
		minCost := cost[i]
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			if cost[v] < minCost {
				minCost = cost[v]
			}
			for _, to := range g[v] {
				if !visited[to] {
					visited[to] = true
					q = append(q, to)
				}
			}
		}
		ans += minCost
	}
	fmt.Fprintln(out, ans)
}
