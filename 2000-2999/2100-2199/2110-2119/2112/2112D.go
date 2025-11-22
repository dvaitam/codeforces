package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v int
}

func dfs(cur, parent int, source bool, g [][]int, ans *[]edge) {
	for _, to := range g[cur] {
		if to == parent {
			continue
		}
		if source {
			*ans = append(*ans, edge{cur, to})
		} else {
			*ans = append(*ans, edge{to, cur})
		}
		dfs(to, cur, !source, g, ans)
	}
}

func solveOne(n int, g [][]int) (bool, []edge) {
	if n == 2 {
		return false, nil
	}

	// find a vertex of degree exactly 2 to be the unique middle
	mid := -1
	for i := 0; i < n; i++ {
		if len(g[i]) == 2 {
			mid = i
			break
		}
	}
	if mid == -1 {
		return false, nil
	}

	a := g[mid][0]
	b := g[mid][1]
	var ans []edge

	// orient a -> mid -> b
	ans = append(ans, edge{a, mid})
	ans = append(ans, edge{mid, b})

	dfs(a, mid, true, g, &ans)  // a is a source
	dfs(b, mid, false, g, &ans) // b is a sink

	return true, ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		ok, ans := solveOne(n, g)
		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}
		fmt.Fprintln(out, "YES")
		for _, e := range ans {
			fmt.Fprintln(out, e.u+1, e.v+1)
		}
	}
}

