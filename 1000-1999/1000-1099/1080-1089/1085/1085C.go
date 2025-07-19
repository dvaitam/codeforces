package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Point holds x and y coordinates
type Point struct {
	x, y int
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pts := make([]Point, 3)
	for i := 0; i < 3; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}
	// sort by x-coordinate
	sort.Slice(pts, func(i, j int) bool {
		return pts[i].x < pts[j].x
	})
	xa, xb, xc := pts[0].x, pts[1].x, pts[2].x
	ya, yb, yc := pts[0].y, pts[1].y, pts[2].y

	var ans []Point
	// horizontal segment from xa to xb at y = ya
	for i := xa; i < xb; i++ {
		ans = append(ans, Point{i, ya})
	}
	// vertical segment at xb from min(ya,yb,yc) to max(ya,yb,yc)
	ymin := min(min(ya, yb), yc)
	ymax := max(max(ya, yb), yc)
	for y := ymin; y <= ymax; y++ {
		ans = append(ans, Point{xb, y})
	}
	// horizontal segment from xb to xc at y = yc
	for i := xc; i > xb; i-- {
		ans = append(ans, Point{i, yc})
	}

	// output result
	fmt.Fprintln(out, len(ans))
	for _, p := range ans {
		fmt.Fprintln(out, p.x, p.y)
	}
}
