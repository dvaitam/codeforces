package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	recipes [][]int
	cost    []int64
	vis     []bool
)

func dfs(u int) int64 {
	if vis[u] {
		return cost[u]
	}
	vis[u] = true
	if len(recipes[u]) > 0 {
		var sum int64
		for _, v := range recipes[u] {
			sum += dfs(v)
		}
		if sum < cost[u] {
			cost[u] = sum
		}
	}
	return cost[u]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		cost = make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &cost[i])
		}
		have := make([]bool, n)
		for i := 0; i < k; i++ {
			var p int
			fmt.Fscan(reader, &p)
			have[p-1] = true
		}
		recipes = make([][]int, n)
		for i := 0; i < n; i++ {
			var m int
			fmt.Fscan(reader, &m)
			if m > 0 {
				recipes[i] = make([]int, m)
				for j := 0; j < m; j++ {
					fmt.Fscan(reader, &recipes[i][j])
					recipes[i][j]--
				}
			}
		}
		vis = make([]bool, n)
		for i := 0; i < n; i++ {
			if have[i] {
				vis[i] = true
				cost[i] = 0
			}
		}
		for i := 0; i < n; i++ {
			dfs(i)
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, cost[i])
		}
		writer.WriteByte('\n')
	}
}
