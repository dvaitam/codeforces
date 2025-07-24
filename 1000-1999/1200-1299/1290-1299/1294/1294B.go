package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct {
	x, y int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pts[i].x, &pts[i].y)
		}
		sort.Slice(pts, func(i, j int) bool {
			if pts[i].x == pts[j].x {
				return pts[i].y < pts[j].y
			}
			return pts[i].x < pts[j].x
		})
		cx, cy := 0, 0
		var path []byte
		ok := true
		for _, p := range pts {
			if p.y < cy {
				ok = false
				break
			}
			for cx < p.x {
				path = append(path, 'R')
				cx++
			}
			for cy < p.y {
				path = append(path, 'U')
				cy++
			}
		}
		if !ok {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, string(path))
		}
	}
}
