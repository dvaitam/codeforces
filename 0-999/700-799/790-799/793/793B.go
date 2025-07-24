package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	x, y  int
	dir   int
	turns int
}

var dx = []int{-1, 0, 1, 0}
var dy = []int{0, 1, 0, -1}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	var sx, sy, tx, ty int
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
		for j := 0; j < m; j++ {
			switch grid[i][j] {
			case 'S':
				sx, sy = i, j
			case 'T':
				tx, ty = i, j
			}
		}
	}

	const INF = 3
	dist := make([][][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([][]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = []int{INF, INF, INF, INF}
		}
	}

	queue := make([]Node, 0)
	for d := 0; d < 4; d++ {
		dist[sx][sy][d] = 0
		queue = append(queue, Node{sx, sy, d, 0})
	}

	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		if cur.turns > dist[cur.x][cur.y][cur.dir] {
			continue
		}
		if cur.turns > 2 {
			continue
		}
		if cur.x == tx && cur.y == ty {
			fmt.Println("YES")
			return
		}
		for nd := 0; nd < 4; nd++ {
			nx := cur.x + dx[nd]
			ny := cur.y + dy[nd]
			nt := cur.turns
			if nd != cur.dir {
				nt++
			}
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if grid[nx][ny] == '*' {
				continue
			}
			if nt < dist[nx][ny][nd] {
				dist[nx][ny][nd] = nt
				queue = append(queue, Node{nx, ny, nd, nt})
			}
		}
	}
	fmt.Println("NO")
}
