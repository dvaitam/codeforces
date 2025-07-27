package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

func bfs(n, m int, w int64, grid []int64, start int) []int64 {
	dist := make([]int64, n*m)
	for i := range dist {
		dist[i] = INF
	}
	queue := make([]int, 0, n*m)
	dist[start] = 0
	queue = append(queue, start)
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		d := dist[v]
		x := v / m
		y := v % m
		for _, dir := range dirs {
			nx := x + dir[0]
			ny := y + dir[1]
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			idx := nx*m + ny
			if grid[idx] == -1 {
				continue
			}
			if dist[idx] == INF {
				dist[idx] = d + w
				queue = append(queue, idx)
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var w int64
	if _, err := fmt.Fscan(reader, &n, &m, &w); err != nil {
		return
	}
	grid := make([]int64, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var val int64
			fmt.Fscan(reader, &val)
			grid[i*m+j] = val
		}
	}

	distStart := bfs(n, m, w, grid, 0)
	distEnd := bfs(n, m, w, grid, n*m-1)

	ans := distStart[n*m-1]
	if ans >= INF {
		ans = INF
	}

	bestStart := INF
	bestEnd := INF

	for idx, val := range grid {
		if val > 0 {
			if distStart[idx] < INF {
				if cost := distStart[idx] + val; cost < bestStart {
					bestStart = cost
				}
			}
			if distEnd[idx] < INF {
				if cost := distEnd[idx] + val; cost < bestEnd {
					bestEnd = cost
				}
			}
		}
	}

	if bestStart < INF && bestEnd < INF {
		if tmp := bestStart + bestEnd; tmp < ans {
			ans = tmp
		}
	}

	if ans >= INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
