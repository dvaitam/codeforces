package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		type point struct {
			x, y int64
		}
		points := make([]point, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &points[i].x, &points[i].y)
		}
		var xs, ys, xt, yt int64
		fmt.Fscan(in, &xs, &ys, &xt, &yt)
		dst := dist2(xs, ys, xt, yt)
		possible := true
		for _, p := range points {
			if dist2(p.x, p.y, xt, yt) <= dst {
				possible = false
				break
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func dist2(x1, y1, x2, y2 int64) int64 {
	dx := x1 - x2
	dy := y1 - y2
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx*dx + dy*dy
}
