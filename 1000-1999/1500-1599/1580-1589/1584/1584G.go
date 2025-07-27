package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y float64
}

func dist2PointSegment(p, a, b Point) float64 {
	abx := b.x - a.x
	aby := b.y - a.y
	apx := p.x - a.x
	apy := p.y - a.y
	ab2 := abx*abx + aby*aby
	if ab2 == 0 {
		dx := apx
		dy := apy
		return dx*dx + dy*dy
	}
	t := (apx*abx + apy*aby) / ab2
	if t < 0 {
		dx := apx
		dy := apy
		return dx*dx + dy*dy
	}
	if t > 1 {
		dx := p.x - b.x
		dy := p.y - b.y
		return dx*dx + dy*dy
	}
	cross := abx*apy - aby*apx
	return (cross * cross) / ab2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var R int
	if _, err := fmt.Fscan(reader, &n, &R); err != nil {
		return
	}
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		points[i] = Point{float64(x), float64(y)}
	}
	r2 := float64(R * R)
	var count int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a := points[i]
			b := points[j]
			ok := true
			for k := 0; k < n; k++ {
				if dist2PointSegment(points[k], a, b) > r2+1e-9 {
					ok = false
					break
				}
			}
			if ok {
				count++
			}
		}
	}
	fmt.Fprintln(writer, count)
}
