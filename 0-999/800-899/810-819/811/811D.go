package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, swapLR, swapUD int
	if _, err := fmt.Fscan(reader, &n, &m, &swapLR, &swapUD); err != nil {
		return
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	dirs := []struct {
		dx, dy int
		c      byte
	}{
		{-1, 0, 'U'},
		{1, 0, 'D'},
		{0, -1, 'L'},
		{0, 1, 'R'},
	}

	parent := make([][]point, n)
	move := make([][]byte, n)
	used := make([][]bool, n)
	for i := range parent {
		parent[i] = make([]point, m)
		move[i] = make([]byte, m)
		used[i] = make([]bool, m)
	}

	q := []point{{0, 0}}
	used[0][0] = true
	var fx, fy int
	found := false
	for len(q) > 0 && !found {
		p := q[0]
		q = q[1:]
		if grid[p.x][p.y] == 'F' {
			fx, fy = p.x, p.y
			found = true
			break
		}
		for _, d := range dirs {
			nx, ny := p.x+d.dx, p.y+d.dy
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if grid[nx][ny] == '*' || used[nx][ny] {
				continue
			}
			used[nx][ny] = true
			parent[nx][ny] = p
			move[nx][ny] = d.c
			q = append(q, point{nx, ny})
		}
	}
	if !found {
		return
	}

	// reconstruct path
	path := []byte{}
	x, y := fx, fy
	for x != 0 || y != 0 {
		path = append(path, move[x][y])
		p := parent[x][y]
		x, y = p.x, p.y
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	for i, c := range path {
		if c == 'L' || c == 'R' {
			if swapLR == 1 {
				if c == 'L' {
					path[i] = 'R'
				} else {
					path[i] = 'L'
				}
			}
		} else {
			if swapUD == 1 {
				if c == 'U' {
					path[i] = 'D'
				} else {
					path[i] = 'U'
				}
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for _, c := range path {
		fmt.Fprintf(writer, "%c\n", c)
	}
}
