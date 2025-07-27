package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct{ x, y int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		x, y := 0, 0
		path := make([]point, len(s))
		for i, ch := range s {
			switch ch {
			case 'L':
				x--
			case 'R':
				x++
			case 'U':
				y++
			case 'D':
				y--
			}
			path[i] = point{x, y}
		}
		ansx, ansy := 0, 0
		for _, p := range path {
			if p.x == 0 && p.y == 0 {
				continue
			}
			cx, cy := 0, 0
			for _, ch := range s {
				nx, ny := cx, cy
				switch ch {
				case 'L':
					nx--
				case 'R':
					nx++
				case 'U':
					ny++
				case 'D':
					ny--
				}
				if nx == p.x && ny == p.y {
					continue
				}
				cx, cy = nx, ny
			}
			if cx == 0 && cy == 0 {
				ansx, ansy = p.x, p.y
				break
			}
		}
		fmt.Fprintln(out, ansx, ansy)
	}
}
