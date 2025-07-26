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
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	g := make([][]int, n)
	gr := make([][]int, n)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		g[x] = append(g[x], y)
		gr[y] = append(gr[y], x)
	}

	// Kosaraju to find SCCs
	order := make([]int, 0, n)
	used := make([]bool, n)
	var dfs1 func(v int)
	dfs1 = func(v int) {
		used[v] = true
		for _, to := range g[v] {
			if !used[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for v := 0; v < n; v++ {
		if !used[v] {
			dfs1(v)
		}
	}

	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	var comps [][]int
	var dfs2 func(v, c int)
	dfs2 = func(v, c int) {
		comp[v] = c
		comps[c] = append(comps[c], v)
		for _, to := range gr[v] {
			if comp[to] == -1 {
				dfs2(to, c)
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			comps = append(comps, []int{})
			dfs2(v, len(comps)-1)
		}
	}

	// Build adjacency lists per component
	type edge struct{ from, to int }
	edges := make([][]edge, len(comps))
	for v := 0; v < n; v++ {
		cv := comp[v]
		for _, to := range g[v] {
			if comp[to] == cv {
				edges[cv] = append(edges[cv], edge{v, to})
			}
		}
	}

	var ansPairs int64
	var ansSingle int64

	// For gcd calculation
	gcd := func(a, b int64) int64 {
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		for b != 0 {
			a, b = b, a%b
		}
		if a < 0 {
			a = -a
		}
		return a
	}

	for cid, nodes := range comps {
		if len(nodes) == 1 && len(edges[cid]) == 0 {
			// isolated node without self-loop
			continue
		}
		// compute gcd and dist for nodes via DFS
		ncnt := len(nodes)
		idMap := make(map[int]int, ncnt)
		for idx, v := range nodes {
			idMap[v] = idx
		}
		dist := make([]int64, ncnt)
		visited := make([]bool, ncnt)
		var gval int64
		var dfs func(idx int)
		dfs = func(idx int) {
			visited[idx] = true
			v := nodes[idx]
			for _, e := range edges[cid] {
				if e.from != v {
					continue
				}
				j := idMap[e.to]
				if !visited[j] {
					dist[j] = dist[idx] + 1
					dfs(j)
				} else {
					diff := dist[idx] + 1 - dist[j]
					gval = gcd(gval, diff)
				}
			}
		}
		for i := range nodes {
			if !visited[i] {
				dfs(i)
			}
		}
		if gval == 0 {
			// no cycles of positive length
			continue
		}
		gmod := int(gval)
		cnt := make([]int64, gmod)
		for i := range nodes {
			r := int(((dist[i] % int64(gmod)) + int64(gmod)) % int64(gmod))
			cnt[r]++
			dist[i] = int64(r)
		}
		km := k % int64(gmod)
		if km == 0 {
			ansSingle += int64(len(nodes))
			for _, c := range cnt {
				ansPairs += c * (c - 1) / 2
			}
		} else if gmod%2 == 0 && km == int64(gmod/2) {
			for r := 0; r < gmod/2; r++ {
				ansPairs += cnt[r] * cnt[r+gmod/2]
			}
		}
	}

	fmt.Fprintln(writer, ansPairs+ansSingle)
}
