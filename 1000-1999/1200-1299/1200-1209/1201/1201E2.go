package main

import (
	"bufio"
	"fmt"
	"os"
)

var dx = []int{1, 1, -1, -1, 2, 2, -2, -2}
var dy = []int{2, -2, 2, -2, 1, -1, 1, -1}

type Point struct{ x, y int }

func bfs(n, m, sx, sy, tx, ty int) []Point {
	const inf = int(1e9)
	dist := make([][]int, n+1)
	px := make([][]int, n+1)
	py := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dist[i] = make([]int, m+1)
		px[i] = make([]int, m+1)
		py[i] = make([]int, m+1)
		for j := 0; j <= m; j++ {
			dist[i][j] = inf
		}
	}
	q := make([]Point, 0)
	dist[sx][sy] = 0
	q = append(q, Point{sx, sy})
	for head := 0; head < len(q); head++ {
		p := q[head]
		if p.x == tx && p.y == ty {
			break
		}
		for k := 0; k < 8; k++ {
			nx := p.x + dx[k]
			ny := p.y + dy[k]
			if nx < 1 || nx > n || ny < 1 || ny > m {
				continue
			}
			if dist[nx][ny] == inf {
				dist[nx][ny] = dist[p.x][p.y] + 1
				px[nx][ny] = p.x
				py[nx][ny] = p.y
				q = append(q, Point{nx, ny})
			}
		}
	}
	// reconstruct
	path := []Point{}
	if dist[tx][ty] == inf {
		return path
	}
	x, y := tx, ty
	for !(x == sx && y == sy) {
		path = append(path, Point{x, y})
		pxOld := px[x][y]
		pyOld := py[x][y]
		x, y = pxOld, pyOld
	}
	path = append(path, Point{sx, sy})
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

	wx, wy := n/2, m/2
	bx, by := n/2+1, m/2

	pathWhite := bfs(n, m, x1, y1, wx, wy)
	pathBlack := bfs(n, m, x2, y2, bx, by)

	chooseWhite := len(pathWhite) <= len(pathBlack)

	if chooseWhite {
		fmt.Fprintln(out, "WHITE")
		out.Flush()
		ox, oy := 0, 0
		for i := 1; i < len(pathWhite); i++ {
			fmt.Fprintf(out, "%d %d\n", pathWhite[i].x, pathWhite[i].y)
			out.Flush()
			fmt.Fscan(in, &ox, &oy)
			if ox == -1 && oy == -1 {
				return
			}
		}
	} else {
		fmt.Fprintln(out, "BLACK")
		out.Flush()
		var ox, oy int
		fmt.Fscan(in, &ox, &oy) // white makes first move
		if ox == -1 && oy == -1 {
			return
		}
		for i := 1; i < len(pathBlack); i++ {
			fmt.Fprintf(out, "%d %d\n", pathBlack[i].x, pathBlack[i].y)
			out.Flush()
			fmt.Fscan(in, &ox, &oy)
			if ox == -1 && oy == -1 {
				return
			}
		}
	}
}
