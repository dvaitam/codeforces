package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func neg(x int) int {
	if x%2 == 0 {
		return x + 1
	}
	return x - 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		grid := make([][]int, 3)
		for i := 0; i < 3; i++ {
			grid[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &grid[i][j])
			}
		}
		N := 2 * n
		g := make([][]int, N)
		gr := make([][]int, N)
		addClause := func(a, b int) {
			na := neg(a)
			nb := neg(b)
			g[na] = append(g[na], b)
			gr[b] = append(gr[b], na)
			g[nb] = append(g[nb], a)
			gr[a] = append(gr[a], nb)
		}
		for j := 0; j < n; j++ {
			var lits [3]int
			for i := 0; i < 3; i++ {
				val := grid[i][j]
				idx := absInt(val) - 1
				if val > 0 {
					lits[i] = 2 * idx
				} else {
					lits[i] = 2*idx + 1
				}
			}
			addClause(lits[0], lits[1])
			addClause(lits[0], lits[2])
			addClause(lits[1], lits[2])
		}

		order := make([]int, 0, N)
		vis := make([]bool, N)
		var dfs1 func(int)
		dfs1 = func(v int) {
			vis[v] = true
			for _, to := range g[v] {
				if !vis[to] {
					dfs1(to)
				}
			}
			order = append(order, v)
		}
		for v := 0; v < N; v++ {
			if !vis[v] {
				dfs1(v)
			}
		}
		comp := make([]int, N)
		for i := range comp {
			comp[i] = -1
		}
		var dfs2 func(int, int)
		dfs2 = func(v, c int) {
			comp[v] = c
			for _, to := range gr[v] {
				if comp[to] == -1 {
					dfs2(to, c)
				}
			}
		}
		cid := 0
		for i := N - 1; i >= 0; i-- {
			v := order[i]
			if comp[v] == -1 {
				dfs2(v, cid)
				cid++
			}
		}
		ok := true
		for i := 0; i < n; i++ {
			if comp[2*i] == comp[2*i+1] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
