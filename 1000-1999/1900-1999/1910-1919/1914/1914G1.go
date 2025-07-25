package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		m := 2 * n
		colors := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &colors[i])
		}

		pos1 := make([]int, n+1)
		pos2 := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pos1[i] = -1
			pos2[i] = -1
		}
		for idx, c := range colors {
			if pos1[c] == -1 {
				pos1[c] = idx
			} else {
				pos2[c] = idx
			}
		}
		for i := 1; i <= n; i++ {
			if pos1[i] > pos2[i] {
				pos1[i], pos2[i] = pos2[i], pos1[i]
			}
		}

		adj := make([][]int, n)
		radj := make([][]int, n)
		for i := 0; i < n; i++ {
			li, ri := pos1[i+1], pos2[i+1]
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				pj1, pj2 := pos1[j+1], pos2[j+1]
				if (li <= pj1 && pj1 <= ri) || (li <= pj2 && pj2 <= ri) {
					adj[i] = append(adj[i], j)
					radj[j] = append(radj[j], i)
				}
			}
		}

		visited := make([]bool, n)
		order := make([]int, 0, n)
		var dfs1 func(int)
		dfs1 = func(v int) {
			visited[v] = true
			for _, u := range adj[v] {
				if !visited[u] {
					dfs1(u)
				}
			}
			order = append(order, v)
		}
		for v := 0; v < n; v++ {
			if !visited[v] {
				dfs1(v)
			}
		}

		comp := make([]int, n)
		for i := range comp {
			comp[i] = -1
		}
		var dfs2 func(int, int)
		dfs2 = func(v, c int) {
			comp[v] = c
			for _, u := range radj[v] {
				if comp[u] == -1 {
					dfs2(u, c)
				}
			}
		}
		cid := 0
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			if comp[v] == -1 {
				dfs2(v, cid)
				cid++
			}
		}

		size := make([]int, cid)
		for _, c := range comp {
			size[c]++
		}
		indeg := make([]int, cid)
		for v := 0; v < n; v++ {
			for _, u := range adj[v] {
				if comp[v] != comp[u] {
					indeg[comp[u]]++
				}
			}
		}

		cntComponents := 0
		ways := int64(1)
		for i := 0; i < cid; i++ {
			if indeg[i] == 0 {
				cntComponents++
				ways = (ways * int64(2*size[i])) % mod
			}
		}

		fmt.Fprintln(writer, cntComponents, ways)
	}
}
