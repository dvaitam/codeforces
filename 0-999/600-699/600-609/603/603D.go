package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct {
	x, y, s int64
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		var a, b, c int64
		fmt.Fscan(in, &a, &b, &c)
		pts[i] = point{x: a * c, y: b * c, s: a*a + b*b}
	}

	var ans int64
	for i := 0; i < n; i++ {
		xi, yi, si := pts[i].x, pts[i].y, pts[i].s
		dup := 0
		slopes := make(map[[2]int64]int)
		for j := i + 1; j < n; j++ {
			xj, yj, sj := pts[j].x, pts[j].y, pts[j].s
			dx := xj*si - xi*sj
			dy := yj*si - yi*sj
			if dx == 0 && dy == 0 {
				dup++
				continue
			}
			g := gcd(abs(dx), abs(dy))
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			slopes[[2]int64{dx, dy}]++
		}
		for _, m := range slopes {
			ans += int64(m * (m - 1) / 2)
			ans += int64(dup * m)
		}
		ans += int64(dup * (dup - 1) / 2)
	}
	fmt.Println(ans)
}
