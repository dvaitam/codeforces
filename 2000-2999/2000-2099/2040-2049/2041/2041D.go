package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	x, y int
	dir  int
	cnt  int
}

var (
	n, m int
	grid []string
)

func encode(x, y, dir, cnt int) int {
	return (((x*m+y)*4)+dir)*3 + (cnt - 1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	grid = make([]string, n)
	var sx, sy, tx, ty int
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grid[i])
		for j := 0; j < m; j++ {
			switch grid[i][j] {
			case 'S':
				sx, sy = i, j
			case 'T':
				tx, ty = i, j
			}
		}
	}

	if sx == tx && sy == ty {
		fmt.Fprintln(out, 0)
		return
	}

	const dirCount = 4
	totalStates := n * m * dirCount * 3
	visited := make([]bool, totalStates)
	dist := make([]int, totalStates)

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	q := make([]node, 0)
	head := 0

	for d := 0; d < dirCount; d++ {
		nx, ny := sx+dx[d], sy+dy[d]
		if nx < 0 || nx >= n || ny < 0 || ny >= m {
			continue
		}
		if grid[nx][ny] == '#' {
			continue
		}
		idx := encode(nx, ny, d, 1)
		if visited[idx] {
			continue
		}
		visited[idx] = true
		dist[idx] = 1
		q = append(q, node{nx, ny, d, 1})
		if nx == tx && ny == ty {
			fmt.Fprintln(out, 1)
			return
		}
	}

	for head < len(q) {
		cur := q[head]
		head++
		curIdx := encode(cur.x, cur.y, cur.dir, cur.cnt)
		for ndir := 0; ndir < dirCount; ndir++ {
			nx, ny := cur.x+dx[ndir], cur.y+dy[ndir]
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if grid[nx][ny] == '#' {
				continue
			}
			ncnt := 1
			if ndir == cur.dir {
				if cur.cnt == 3 {
					continue
				}
				ncnt = cur.cnt + 1
			}
			idx := encode(nx, ny, ndir, ncnt)
			if visited[idx] {
				continue
			}
			visited[idx] = true
			dist[idx] = dist[curIdx] + 1
			if nx == tx && ny == ty {
				fmt.Fprintln(out, dist[idx])
				return
			}
			q = append(q, node{nx, ny, ndir, ncnt})
		}
	}

	fmt.Fprintln(out, -1)
}
