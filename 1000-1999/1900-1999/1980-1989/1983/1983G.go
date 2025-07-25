package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

var (
	g     [][]int
	up    [][]int
	depth []int
	val   []int
)

func dfs(u, p int) {
	up[u][0] = p
	for i := 1; i < LOG; i++ {
		up[u][i] = up[up[u][i-1]][i-1]
	}
	for _, v := range g[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs(v, u)
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			a = up[a][i]
		}
	}
	if a == b {
		return a
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[a][i] != up[b][i] {
			a = up[a][i]
			b = up[b][i]
		}
	}
	return up[a][0]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		g = make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		val = make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &val[i])
		}
		up = make([][]int, n+1)
		depth = make([]int, n+1)
		for i := 0; i <= n; i++ {
			up[i] = make([]int, LOG)
		}
		depth[1] = 0
		dfs(1, 1)

		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			l := lca(x, y)
			ans := int64(0)
			idx := 0
			u := x
			for u != l {
				ans += int64(val[u] ^ idx)
				idx++
				u = up[u][0]
			}
			ans += int64(val[l] ^ idx)
			idx++
			var stack []int
			v := y
			for v != l {
				stack = append(stack, v)
				v = up[v][0]
			}
			for i := len(stack) - 1; i >= 0; i-- {
				ans += int64(val[stack[i]] ^ idx)
				idx++
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
