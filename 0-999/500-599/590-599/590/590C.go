package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

var dirs = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	const INF = int(1e9)
	dist := make([][]int, 3)
	for k := 0; k < 3; k++ {
		dist[k] = make([]int, n*m)
		for i := range dist[k] {
			dist[k][i] = INF
		}
	}

	for k := 0; k < 3; k++ {
		dq := list.New()
		// initialize with all cells of state k+1
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == byte('1'+k) {
					idx := i*m + j
					dist[k][idx] = 0
					dq.PushBack(idx)
				}
			}
		}
		for dq.Len() > 0 {
			e := dq.Front()
			dq.Remove(e)
			idx := e.Value.(int)
			x, y := idx/m, idx%m
			dcur := dist[k][idx]
			for _, d := range dirs {
				nx, ny := x+d[0], y+d[1]
				if nx < 0 || nx >= n || ny < 0 || ny >= m {
					continue
				}
				if grid[nx][ny] == '#' {
					continue
				}
				add := 0
				if grid[nx][ny] == '.' {
					add = 1
				}
				nidx := nx*m + ny
				nd := dcur + add
				if nd < dist[k][nidx] {
					dist[k][nidx] = nd
					if add == 1 {
						dq.PushBack(nidx)
					} else {
						dq.PushFront(nidx)
					}
				}
			}
		}
	}

	best := INF
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			idx := i*m + j
			d1, d2, d3 := dist[0][idx], dist[1][idx], dist[2][idx]
			if d1 == INF || d2 == INF || d3 == INF {
				continue
			}
			total := d1 + d2 + d3
			if grid[i][j] == '.' {
				total -= 2
			}
			if total < best {
				best = total
			}
		}
	}

	if best >= INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, best)
	}
}
