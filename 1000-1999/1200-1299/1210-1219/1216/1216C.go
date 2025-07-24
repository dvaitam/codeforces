package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rect struct {
	x1, y1, x2, y2 int
}

func intersect(a, b Rect) Rect {
	x1 := max(a.x1, b.x1)
	y1 := max(a.y1, b.y1)
	x2 := min(a.x2, b.x2)
	y2 := min(a.y2, b.y2)
	if x1 >= x2 || y1 >= y2 {
		return Rect{x1, y1, x1, y1}
	}
	return Rect{x1, y1, x2, y2}
}

func area(r Rect) int64 {
	if r.x2 <= r.x1 || r.y2 <= r.y1 {
		return 0
	}
	return int64(r.x2-r.x1) * int64(r.y2-r.y1)
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

	var w, b1, b2 Rect
	fmt.Fscan(in, &w.x1, &w.y1, &w.x2, &w.y2)
	fmt.Fscan(in, &b1.x1, &b1.y1, &b1.x2, &b1.y2)
	fmt.Fscan(in, &b2.x1, &b2.y1, &b2.x2, &b2.y2)

	wArea := area(w)
	cov1 := area(intersect(w, b1))
	cov2 := area(intersect(w, b2))
	covBoth := area(intersect(intersect(w, b1), b2))
	covered := cov1 + cov2 - covBoth
	if covered < wArea {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
