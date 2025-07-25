package main

import (
	"bufio"
	"fmt"
	"os"
)

var dr = []int{-1, -1, -1, 0, 0, 1, 1, 1}
var dc = []int{-1, 0, 1, -1, 1, -1, 0, 1}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	lines := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		lines[i] = []byte(s)
	}

	n2 := n + 2
	m2 := m + 2
	size := n2 * m2

	grid := make([]byte, size)
	dist := make([]int, size)
	queue := make([]int, 0, size)

	for i := 0; i < n2; i++ {
		for j := 0; j < m2; j++ {
			id := i*m2 + j
			var ch byte
			if i == 0 || j == 0 || i == n+1 || j == m+1 {
				ch = '.'
			} else {
				ch = lines[i-1][j-1]
			}
			grid[id] = ch
			if ch == '.' {
				dist[id] = 0
				queue = append(queue, id)
			} else {
				dist[id] = -1
			}
		}
	}

	// BFS to compute distance from each cell to nearest '.' (Chebyshev metric)
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		d := dist[id]
		r := id / m2
		c := id % m2
		for k := 0; k < 8; k++ {
			nr := r + dr[k]
			nc := c + dc[k]
			if nr >= 0 && nr < n2 && nc >= 0 && nc < m2 {
				nid := nr*m2 + nc
				if dist[nid] == -1 {
					dist[nid] = d + 1
					queue = append(queue, nid)
				}
			}
		}
	}

	maxd := 0
	for i := 0; i < size; i++ {
		if grid[i] == 'X' && dist[i] > maxd {
			maxd = dist[i]
		}
	}
	if maxd > 0 {
		maxd--
	}

	check := func(mid int) bool {
		visited := make([]bool, size)
		q := make([]int, 0)
		qs := make([]int, 0)
		for i := 0; i < size; i++ {
			if dist[i] > mid {
				visited[i] = true
				q = append(q, i)
				qs = append(qs, 0)
			}
		}
		for head := 0; head < len(q); head++ {
			id := q[head]
			step := qs[head]
			if step == mid {
				continue
			}
			r := id / m2
			c := id % m2
			for k := 0; k < 8; k++ {
				nr := r + dr[k]
				nc := c + dc[k]
				if nr >= 0 && nr < n2 && nc >= 0 && nc < m2 {
					nid := nr*m2 + nc
					if !visited[nid] {
						visited[nid] = true
						q = append(q, nid)
						qs = append(qs, step+1)
					}
				}
			}
		}
		for i := 0; i < size; i++ {
			r := i / m2
			c := i % m2
			if r == 0 || c == 0 || r == n+1 || c == m+1 {
				if visited[i] {
					return false
				}
			} else {
				if visited[i] != (grid[i] == 'X') {
					return false
				}
			}
		}
		return true
	}

	lo, hi := 0, maxd
	best := 0
	for lo <= hi {
		mid := (lo + hi) / 2
		if check(mid) {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	// reconstruct initial burning cells
	ans := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			id := (i+1)*m2 + (j + 1)
			if dist[id] > best {
				row[j] = 'X'
			} else {
				row[j] = '.'
			}
		}
		ans[i] = row
	}

	fmt.Fprintln(writer, best)
	for i := 0; i < n; i++ {
		writer.Write(ans[i])
		writer.WriteByte('\n')
	}
}
