package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pawn struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	pawns := make([]Pawn, n)
	const MAX = 500000
	cols := make([][]int, MAX+2)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		pawns[i] = Pawn{x, y}
		cols[y] = append(cols[y], x)
	}
	for i := range cols {
		if len(cols[i]) > 0 {
			sort.Ints(cols[i])
		}
	}
	id := make(map[[2]int]int, n)
	for i, p := range pawns {
		id[[2]int{p.x, p.y}] = i
	}

	adj := make([][]int, n)
	const INF = int(1e9)
	for y, xs := range cols {
		if len(xs) == 0 {
			continue
		}
		for idx, x := range xs {
			next := INF
			if idx+1 < len(xs) {
				next = xs[idx+1]
			}
			u := id[[2]int{x, y}]
			for _, d := range [2]int{-1, 1} {
				y2 := y + d
				if y2 < 1 || y2 > MAX {
					continue
				}
				xsAdj := cols[y2]
				if len(xsAdj) == 0 {
					continue
				}
				j := sort.SearchInts(xsAdj, x+1)
				if j < len(xsAdj) {
					x2 := xsAdj[j]
					if x2 < next {
						v := id[[2]int{x2, y2}]
						adj[u] = append(adj[u], v)
					}
				}
			}
		}
	}

	pairU := make([]int, n)
	pairV := make([]int, n)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		pairU[i] = -1
		pairV[i] = -1
	}

	queue := make([]int, 0, n)
	bfs := func() bool {
		queue = queue[:0]
		found := false
		for u := 0; u < n; u++ {
			if pairU[u] == -1 {
				dist[u] = 0
				queue = append(queue, u)
			} else {
				dist[u] = -1
			}
		}
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				pu := pairV[v]
				if pu != -1 && dist[pu] == -1 {
					dist[pu] = dist[u] + 1
					queue = append(queue, pu)
				}
				if pu == -1 {
					found = true
				}
			}
		}
		return found
	}

	var dfs func(int) bool
	dfs = func(u int) bool {
		for _, v := range adj[u] {
			pu := pairV[v]
			if pu == -1 || (dist[pu] == dist[u]+1 && dfs(pu)) {
				pairU[u] = v
				pairV[v] = u
				return true
			}
		}
		dist[u] = -1
		return false
	}

	matching := 0
	for bfs() {
		for u := 0; u < n; u++ {
			if pairU[u] == -1 && dfs(u) {
				matching++
			}
		}
	}

	result := n - matching
	fmt.Fprintln(writer, result)
}
