package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	parent := make([]int, n+1)
	dist := make([]int64, n+1)
	switchInit := make([]int, n+1)

	for i := 0; i < n-1; i++ {
		var u, v int
		var d int64
		fmt.Fscan(in, &u, &v, &d)
		g[u] = append(g[u], Edge{v, d})
		parent[v] = u
		dist[v] = dist[u] + d
		switchInit[u] = v // last outgoing edge in input
	}

	type Train struct {
		s int
		t int64
	}
	trains := make([]Train, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &trains[i].s, &trains[i].t)
	}

	// maps time -> used
	used := make(map[int64]bool)
	orientation := make([]int, n+1)
	lastUsed := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		orientation[i] = switchInit[i]
	}

	var explode int64 = -1
	var changes int

	for _, tr := range trains {
		// build path from root to destination
		pathNodes := make([]int, 0)
		cur := tr.s
		for cur != 0 {
			pathNodes = append(pathNodes, cur)
			cur = parent[cur]
		}
		// reverse to root->...
		for i, j := 0, len(pathNodes)-1; i < j; i, j = i+1, j-1 {
			pathNodes[i], pathNodes[j] = pathNodes[j], pathNodes[i]
		}
		time := tr.t
		for i := 0; i+1 < len(pathNodes); i++ {
			u := pathNodes[i]
			v := pathNodes[i+1]
			// find weight
			var w int64
			for _, e := range g[u] {
				if e.to == v {
					w = e.w
					break
				}
			}
			if orientation[u] != v {
				if used[time] {
					if explode == -1 || time < explode {
						explode = time
					}
				} else {
					used[time] = true
					orientation[u] = v
					changes++
				}
			}
			lastUsed[u] = time
			time += w
		}
	}

	fmt.Fprintf(out, "%d %d\n", explode, changes)
}
