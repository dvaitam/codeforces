package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

var (
	g       [][]int
	nearest []int
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dfs(v, p, dist int) int {
	children := true
	sum := 0
	for _, u := range g[v] {
		if u == p {
			continue
		}
		c := dfs(u, v, dist+1)
		if c < 0 {
			children = false
		}
		if nearest[u]+1 < nearest[v] {
			nearest[v] = nearest[u] + 1
		}
		sum += c
	}
	if nearest[v] <= dist {
		return 1
	}
	if sum == 0 || !children {
		return -1
	}
	return sum
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		g = make([][]int, n)
		nearest = make([]int, n)
		for i := range nearest {
			nearest[i] = n + 5
		}
		for i := 0; i < k; i++ {
			var x int
			fmt.Fscan(reader, &x)
			x--
			nearest[x] = 0
		}
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		ans := dfs(0, -1, 0)
		fmt.Fprintln(writer, ans)
	}
}
