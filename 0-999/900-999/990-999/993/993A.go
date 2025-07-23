package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x int
	y int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func orientation(a, b, c Point) int64 {
	return int64(b.x-a.x)*int64(c.y-a.y) - int64(b.y-a.y)*int64(c.x-a.x)
}

func onSegment(a, b, c Point) bool { // b is point, segment ac
	if b.x >= min(a.x, c.x) && b.x <= max(a.x, c.x) &&
		b.y >= min(a.y, c.y) && b.y <= max(a.y, c.y) {
		return orientation(a, c, b) == 0
	}
	return false
}

func segmentsIntersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if (o1 > 0 && o2 < 0 || o1 < 0 && o2 > 0) &&
		(o3 > 0 && o4 < 0 || o3 < 0 && o4 > 0) {
		return true
	}
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}
	return false
}

func pointInside(p Point, poly []Point) bool {
	n := len(poly)
	sign := 0
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		val := orientation(a, b, p)
		if val == 0 {
			if onSegment(a, p, b) {
				return true
			}
			continue
		}
		if sign == 0 {
			if val > 0 {
				sign = 1
			} else {
				sign = -1
			}
		} else {
			if (val > 0 && sign < 0) || (val < 0 && sign > 0) {
				return false
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	points1 := make([]Point, 4)
	points2 := make([]Point, 4)
	for i := 0; i < 4; i++ {
		if _, err := fmt.Fscan(in, &points1[i].x, &points1[i].y); err != nil {
			return
		}
	}
	for i := 0; i < 4; i++ {
		if _, err := fmt.Fscan(in, &points2[i].x, &points2[i].y); err != nil {
			return
		}
	}

	// edges of squares
	for i := 0; i < 4; i++ {
		if pointInside(points1[i], points2) {
			fmt.Println("Yes")
			return
		}
		if pointInside(points2[i], points1) {
			fmt.Println("Yes")
			return
		}
	}
	// check edge intersection
	for i := 0; i < 4; i++ {
		a1 := points1[i]
		a2 := points1[(i+1)%4]
		for j := 0; j < 4; j++ {
			b1 := points2[j]
			b2 := points2[(j+1)%4]
			if segmentsIntersect(a1, a2, b1, b2) {
				fmt.Println("Yes")
				return
			}
		}
	}
	fmt.Println("No")
}
