package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x int64
	y int64
}

func onLine(a, b, c Point) bool {
	return (b.x-a.x)*(c.y-a.y) == (b.y-a.y)*(c.x-a.x)
}

func check(points []Point, i, j int) bool {
	a := points[i]
	b := points[j]
	var others []Point
	for k := 0; k < len(points); k++ {
		if !onLine(a, b, points[k]) {
			others = append(others, points[k])
		}
	}
	if len(others) <= 2 {
		return true
	}
	c := others[0]
	d := others[1]
	for _, p := range others[2:] {
		if !onLine(c, d, p) {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		var x, y int64
		fmt.Fscan(reader, &x, &y)
		points[i] = Point{x, y}
	}
	if n <= 4 {
		fmt.Println("YES")
		return
	}
	pairs := [][2]int{{0, 1}, {0, 2}, {1, 2}}
	for _, pr := range pairs {
		if check(points, pr[0], pr[1]) {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
