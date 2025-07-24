package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type HopcroftKarp struct {
	n     int
	m     int
	adj   [][]int
	pairU []int
	pairV []int
	dist  []int
}

func NewHopcroftKarp(n, m int, adj [][]int) *HopcroftKarp {
	hk := &HopcroftKarp{n: n, m: m, adj: adj}
	hk.pairU = make([]int, n)
	hk.pairV = make([]int, m)
	hk.dist = make([]int, n)
	for i := 0; i < n; i++ {
		hk.pairU[i] = -1
	}
	for j := 0; j < m; j++ {
		hk.pairV[j] = -1
	}
	return hk
}

func (hk *HopcroftKarp) bfs() bool {
	queue := make([]int, 0)
	for i := 0; i < hk.n; i++ {
		if hk.pairU[i] == -1 {
			hk.dist[i] = 0
			queue = append(queue, i)
		} else {
			hk.dist[i] = -1
		}
	}
	found := false
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range hk.adj[u] {
			pu := hk.pairV[v]
			if pu != -1 && hk.dist[pu] == -1 {
				hk.dist[pu] = hk.dist[u] + 1
				queue = append(queue, pu)
			} else if pu == -1 {
				found = true
			}
		}
	}
	return found
}

func (hk *HopcroftKarp) dfs(u int) bool {
	for _, v := range hk.adj[u] {
		pu := hk.pairV[v]
		if pu == -1 || (hk.dist[pu] == hk.dist[u]+1 && hk.dfs(pu)) {
			hk.pairU[u] = v
			hk.pairV[v] = u
			return true
		}
	}
	hk.dist[u] = -1
	return false
}

func (hk *HopcroftKarp) MaxMatching() int {
	matching := 0
	for hk.bfs() {
		for i := 0; i < hk.n; i++ {
			if hk.pairU[i] == -1 && hk.dfs(i) {
				matching++
			}
		}
	}
	return matching
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &q)
	board := make([][]bool, n)
	for i := 0; i < n; i++ {
		row := make([]bool, n)
		for j := 0; j < n; j++ {
			row[j] = true
		}
		board[i] = row
	}
	for k := 0; k < q; k++ {
		var x1, y1, x2, y2 int
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		for i := x1 - 1; i <= x2-1; i++ {
			for j := y1 - 1; j <= y2-1; j++ {
				board[i][j] = false
			}
		}
	}
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		cols := make([]int, 0)
		for j := 0; j < n; j++ {
			if board[i][j] {
				cols = append(cols, j)
			}
		}
		adj[i] = cols
	}
	hk := NewHopcroftKarp(n, n, adj)
	ans := hk.MaxMatching()
	fmt.Println(ans)
}
