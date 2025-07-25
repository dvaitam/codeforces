package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

type Point struct {
	x int
	y int
}

type Sensor struct {
	x float64
	y float64
	r float64
}

func dist2PointSegment(px, py, x1, y1, x2, y2 float64) float64 {
	vx := x2 - x1
	vy := y2 - y1
	ux := px - x1
	uy := py - y1
	l2 := vx*vx + vy*vy
	t := (vx*ux + vy*uy) / l2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	dx := x1 + t*vx - px
	dy := y1 + t*vy - py
	return dx*dx + dy*dy
}

func clearSegment(x1, y1, x2, y2 float64, sensors []Sensor) bool {
	for _, s := range sensors {
		d2 := dist2PointSegment(s.x, s.y, x1, y1, x2, y2)
		if d2+1e-8 < s.r*s.r {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var N int
	var xL, yL, xR, yR int
	if _, err := fmt.Fscan(in, &N, &xL, &yL, &xR, &yR); err != nil {
		return
	}
	var xs, ys int
	fmt.Fscan(in, &xs, &ys)
	var xt, yt int
	fmt.Fscan(in, &xt, &yt)

	sensors := make([]Sensor, N)
	for i := 0; i < N; i++ {
		var xi, yi, ri int
		fmt.Fscan(in, &xi, &yi, &ri)
		sensors[i] = Sensor{float64(xi), float64(yi), float64(ri)}
	}

	if xs == xt && ys == yt {
		fmt.Println(0)
		return
	}

	dxs := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dys := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	width := xR - xL + 1
	height := yR - yL + 1

	visited := make([][]bool, width)
	prev := make([][]Point, width)
	for i := range visited {
		visited[i] = make([]bool, height)
		prev[i] = make([]Point, height)
	}

	q := list.New()
	start := Point{xs, ys}
	q.PushBack(start)
	visited[xs-xL][ys-yL] = true

	found := false
	for q.Len() > 0 && !found {
		e := q.Front()
		q.Remove(e)
		p := e.Value.(Point)
		if p.x == xt && p.y == yt {
			found = true
			break
		}
		for dir := 0; dir < 8; dir++ {
			nx := p.x + dxs[dir]
			ny := p.y + dys[dir]
			if nx < xL || nx > xR || ny < yL || ny > yR {
				continue
			}
			if visited[nx-xL][ny-yL] {
				continue
			}
			if !clearSegment(float64(p.x), float64(p.y), float64(nx), float64(ny), sensors) {
				continue
			}
			visited[nx-xL][ny-yL] = true
			prev[nx-xL][ny-yL] = p
			q.PushBack(Point{nx, ny})
		}
	}

	if !visited[xt-xL][yt-yL] {
		fmt.Println(0)
		return
	}

	var path []Point
	cur := Point{xt, yt}
	for {
		path = append(path, cur)
		if cur.x == xs && cur.y == ys {
			break
		}
		cur = prev[cur.x-xL][cur.y-yL]
	}

	// reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// smooth path
	var smooth []Point
	i := 0
	for i < len(path)-1 {
		j := len(path) - 1
		for j > i+1 {
			if clearSegment(float64(path[i].x), float64(path[i].y), float64(path[j].x), float64(path[j].y), sensors) {
				break
			}
			j--
		}
		smooth = append(smooth, path[i])
		i = j
	}
	smooth = append(smooth, path[len(path)-1])

	// output
	m := len(smooth) - 2
	if m < 0 {
		m = 0
	}
	fmt.Println(m)
	for k := 1; k <= m; k++ {
		fmt.Printf("%d %d\n", smooth[k].x, smooth[k].y)
	}
}
