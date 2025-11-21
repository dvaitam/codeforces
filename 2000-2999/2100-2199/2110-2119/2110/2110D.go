package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	w  int64
}

const inf int64 = 1<<62 - 1

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}

		adj := make([][]edge, n+1)
		for i := 0; i < m; i++ {
			var s, t int
			var w int64
			fmt.Fscan(in, &s, &t, &w)
			adj[s] = append(adj[s], edge{to: t, w: w})
		}

		need := make([]int64, n+1)
		final := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			need[i] = inf
			final[i] = inf
		}
		need[n] = 0
		final[n] = 0

		for v := n - 1; v >= 1; v-- {
			bestEntry := inf
			bestFinal := inf
			for _, e := range adj[v] {
				to := e.to
				if need[to] == inf {
					continue
				}
				L := need[to]
				if L < e.w {
					L = e.w
				}
				entry := L - b[v]
				if entry < 0 {
					entry = 0
				}
				finalCand := final[to]
				if finalCand < L {
					finalCand = L
				}
				if entry < bestEntry || (entry == bestEntry && finalCand < bestFinal) {
					bestEntry = entry
					bestFinal = finalCand
				}
			}
			need[v] = bestEntry
			final[v] = bestFinal
		}

		if need[1] != 0 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, final[1])
		}
	}
}
