package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}

	color := make([]int8, n)
	for i := range color {
		color[i] = -1
	}

	var t int
	var ways int64

	var queue []int
	var bipartitePairs int64
	for v := 0; v < n; v++ {
		if color[v] != -1 {
			continue
		}
		// BFS to check bipartiteness
		queue = queue[:0]
		queue = append(queue, v)
		color[v] = 0
		cnt := [2]int64{}
		cnt[0]++
		isBip := true
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, to := range g[u] {
				if color[to] == -1 {
					color[to] = color[u] ^ 1
					cnt[color[to]]++
					queue = append(queue, to)
				} else if color[to] == color[u] {
					isBip = false
				}
			}
		}
		if !isBip {
			fmt.Println(0, 1)
			return
		}
		// add number of pairs within same color
		bipartitePairs += cnt[0]*(cnt[0]-1)/2 + cnt[1]*(cnt[1]-1)/2
	}

	if bipartitePairs > 0 {
		t = 1
		ways = bipartitePairs
	} else if m > 0 {
		t = 2
		ways = int64(m) * int64(n-2)
	} else {
		t = 3
		ways = int64(n) * int64(n-1) * int64(n-2) / 6
	}
	fmt.Println(t, ways)
}
