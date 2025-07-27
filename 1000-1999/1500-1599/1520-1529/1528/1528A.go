package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	g        [][]int
	l, r     []int64
	dp0, dp1 []int64
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func dfs(u, p int) {
	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u)
		// If u uses l[u]
		option0 := absInt64(l[u]-l[v]) + dp0[v]
		if val := absInt64(l[u]-r[v]) + dp1[v]; val > option0 {
			option0 = val
		}
		dp0[u] += option0
		// If u uses r[u]
		option1 := absInt64(r[u]-l[v]) + dp0[v]
		if val := absInt64(r[u]-r[v]) + dp1[v]; val > option1 {
			option1 = val
		}
		dp1[u] += option1
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		l = make([]int64, n+1)
		r = make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
		}
		g = make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		dp0 = make([]int64, n+1)
		dp1 = make([]int64, n+1)
		dfs(1, 0)
		if dp0[1] > dp1[1] {
			fmt.Fprintln(writer, dp0[1])
		} else {
			fmt.Fprintln(writer, dp1[1])
		}
	}
}
