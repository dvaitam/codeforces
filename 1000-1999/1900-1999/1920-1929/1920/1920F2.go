// Solution for problem F2 (see problemF2.txt).
// This version uses the same approach as F1 but handles larger q.

package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, m, q int
var grid [][]byte
var dist []int
var islands []int

func idx(r, c int) int { return r*m + c }

func inBounds(r, c int) bool { return r >= 0 && r < n && c >= 0 && c < m }

// BFS from all volcano cells to compute Manhattan distance to nearest volcano
func computeDist() {
	dist = make([]int, n*m)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'v' {
				id := idx(i, j)
				dist[id] = 0
				queue = append(queue, id)
			}
			if grid[i][j] == '#' {
				islands = append(islands, idx(i, j))
			}
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		d := dist[v] + 1
		if r > 0 {
			u := v - m
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if r+1 < n {
			u := v + m
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if c > 0 {
			u := v - 1
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if c+1 < m {
			u := v + 1
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
	}
}

// Build 4-neigh component reachable from start with dist >= t on non-island cells
func buildComponent(start, t int) []bool {
	comp := make([]bool, n*m)
	if dist[start] < t {
		return comp
	}
	queue := []int{start}
	comp[start] = true
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		if r > 0 {
			u := v - m
			if !comp[u] && grid[r-1][c] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if r+1 < n {
			u := v + m
			if !comp[u] && grid[r+1][c] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if c > 0 {
			u := v - 1
			if !comp[u] && grid[r][c-1] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if c+1 < m {
			u := v + 1
			if !comp[u] && grid[r][c+1] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
	}
	return comp
}

// Check if island can reach border without crossing cells in comp
func islandToBorder(comp []bool) bool {
	visited := make([]bool, n*m)
	queue := make([]int, len(islands))
	copy(queue, islands)
	for _, id := range islands {
		visited[id] = true
	}
	dirs := [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		if r == 0 || r == n-1 || c == 0 || c == m-1 {
			return true
		}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if !inBounds(nr, nc) {
				continue
			}
			ni := idx(nr, nc)
			if comp[ni] || visited[ni] {
				continue
			}
			visited[ni] = true
			queue = append(queue, ni)
		}
	}
	return false
}

func check(start, t int) bool {
	if dist[start] < t {
		return false
	}
	comp := buildComponent(start, t)
	if islandToBorder(comp) {
		return false
	}
	return true
}

func solveQuery(x, y int) int {
	start := idx(x-1, y-1)
	left, right := 0, dist[start]
	ans := 0
	for left <= right {
		mid := (left + right) / 2
		if check(start, mid) {
			ans = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	grid = make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}
	computeDist()
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		fmt.Fprintln(out, solveQuery(x, y))
	}
}
