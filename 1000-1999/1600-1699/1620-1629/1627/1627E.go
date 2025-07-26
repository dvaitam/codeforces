package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for problemE.txt (No Escape).
// We model only rooms involved in ladders plus (1,1) and (n,m).
// For each floor we relax distances horizontally and propagate through ladders.

type nodeKey struct {
	f int
	c int
}

type edge struct {
	to int
	w  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		x := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &x[i])
		}

		type ladder struct {
			a, b, c, d int
			h          int64
			s, e       int
		}
		ladders := make([]ladder, k)

		idxMap := make(map[nodeKey]int)
		floors := []int{}
		cols := []int{}
		nodesByFloor := make(map[int][]int)
		var getIndex func(int, int) int
		getIndex = func(f, c int) int {
			key := nodeKey{f, c}
			if id, ok := idxMap[key]; ok {
				return id
			}
			id := len(floors)
			idxMap[key] = id
			floors = append(floors, f)
			cols = append(cols, c)
			nodesByFloor[f] = append(nodesByFloor[f], id)
			return id
		}

		startIdx := getIndex(1, 1)
		endIdx := getIndex(n, m)

		for i := 0; i < k; i++ {
			var a, b, c, d int
			var h int64
			fmt.Fscan(in, &a, &b, &c, &d, &h)
			s := getIndex(a, b)
			e := getIndex(c, d)
			ladders[i] = ladder{a, b, c, d, h, s, e}
		}

		for f, ids := range nodesByFloor {
			sort.Slice(ids, func(i, j int) bool { return cols[ids[i]] < cols[ids[j]] })
			nodesByFloor[f] = ids
		}

		edges := make([][]edge, len(floors))
		for _, ld := range ladders {
			edges[ld.s] = append(edges[ld.s], edge{ld.e, -ld.h})
		}

		const INF int64 = 1 << 60
		dist := make([]int64, len(floors))
		for i := range dist {
			dist[i] = INF
		}
		dist[startIdx] = 0

		for f := 1; f <= n; f++ {
			ids := nodesByFloor[f]
			if len(ids) == 0 {
				continue
			}
			// left to right
			for i := 1; i < len(ids); i++ {
				prev := ids[i-1]
				cur := ids[i]
				cost := int64(cols[cur]-cols[prev]) * x[f]
				if dist[prev]+cost < dist[cur] {
					dist[cur] = dist[prev] + cost
				}
			}
			// right to left
			for i := len(ids) - 2; i >= 0; i-- {
				next := ids[i+1]
				cur := ids[i]
				cost := int64(cols[next]-cols[cur]) * x[f]
				if dist[next]+cost < dist[cur] {
					dist[cur] = dist[next] + cost
				}
			}
			// use ladders
			for _, id := range ids {
				if dist[id] == INF {
					continue
				}
				for _, e := range edges[id] {
					if dist[id]+e.w < dist[e.to] {
						dist[e.to] = dist[id] + e.w
					}
				}
			}
		}

		if dist[endIdx] >= INF/2 {
			fmt.Fprintln(out, "NO ESCAPE")
		} else {
			fmt.Fprintln(out, dist[endIdx])
		}
	}
}
