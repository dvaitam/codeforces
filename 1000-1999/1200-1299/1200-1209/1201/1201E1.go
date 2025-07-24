package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

var dirs = []Point{{1, 2}, {1, -2}, {-1, 2}, {-1, -2}, {2, 1}, {2, -1}, {-2, 1}, {-2, -1}}

func bfs(n, m int, start, target Point) []Point {
	visited := make([][]bool, n+1)
	parent := make([][]Point, n+1)
	for i := range visited {
		visited[i] = make([]bool, m+1)
		parent[i] = make([]Point, m+1)
	}
	q := []Point{start}
	visited[start.x][start.y] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur == target {
			break
		}
		for _, d := range dirs {
			nx, ny := cur.x+d.x, cur.y+d.y
			if nx < 1 || nx > n || ny < 1 || ny > m {
				continue
			}
			if !visited[nx][ny] {
				visited[nx][ny] = true
				parent[nx][ny] = cur
				q = append(q, Point{nx, ny})
			}
		}
	}
	// reconstruct
	path := []Point{}
	cur := target
	for cur != start {
		path = append(path, cur)
		cur = parent[cur.x][cur.y]
	}
	path = append(path, start)
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var x1, y1, x2, y2 int
	fmt.Fscan(in, &x1, &y1, &x2, &y2)

	white := Point{x1, y1}
	target := Point{n / 2, m / 2}

	path := bfs(n, m, white, target)

	fmt.Fprintln(out, "WHITE")
	out.Flush()

	var bx, by int
	// we start from path[1]
	for i := 1; i < len(path); i++ {
		fmt.Fprintf(out, "%d %d\n", path[i].x, path[i].y)
		out.Flush()
		if _, err := fmt.Fscan(in, &bx, &by); err != nil {
			return
		}
		if bx == -1 && by == -1 {
			return
		}
	}
}
