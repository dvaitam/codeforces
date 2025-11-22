package main

import (
	"bufio"
	"fmt"
	"os"
)

// Checks whether two vertices u, v in a graph are twins:
// they have identical adjacency to all other vertices (their mutual edge is irrelevant).
func isTwin(u, v int, adj [][]uint64) bool {
	wu, bu := u>>6, uint(u&63)
	wv, bv := v>>6, uint(v&63)
	for i := 0; i < len(adj[u]); i++ {
		a := adj[u][i]
		b := adj[v][i]
		if i == wu {
			a &^= 1 << bu
		}
		if i == wv {
			a &^= 1 << bv
		}
		if i == wu {
			b &^= 1 << bu
		}
		if i == wv {
			b &^= 1 << bv
		}
		if a != b {
			return false
		}
	}
	return true
}

// A graph is a possible superb graph iff it has no twin vertices.
func twinFree(n int, edges [][2]int) bool {
	words := (n + 63) >> 6
	adj := make([][]uint64, n)
	for i := range adj {
		adj[i] = make([]uint64, words)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u][v>>6] |= 1 << uint(v&63)
		adj[v][u>>6] |= 1 << uint(u&63)
	}

	for u := 0; u < n; u++ {
		for v := u + 1; v < n; v++ {
			if isTwin(u, v, adj) {
				return false
			}
		}
	}
	return true
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		ok := true
		for i := 0; i < k; i++ {
			var m int
			fmt.Fscan(in, &m)
			edges := make([][2]int, m)
			for j := 0; j < m; j++ {
				var u, v int
				fmt.Fscan(in, &u, &v)
				u--
				v--
				edges[j] = [2]int{u, v}
			}
			if ok && !twinFree(n, edges) {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

