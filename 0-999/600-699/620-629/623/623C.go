package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	// helper to compute squared diameter given assignment arrays
	compute := func(assignX []bool) int64 {
		var xMin, xMax, yMin, yMax int64
		var ax, ay int64
		firstX := true
		firstY := true
		for i := 0; i < n; i++ {
			if assignX[i] {
				if firstX {
					xMin, xMax = xs[i], xs[i]
					ax = abs64(xs[i])
					firstX = false
				} else {
					if xs[i] < xMin {
						xMin = xs[i]
					}
					if xs[i] > xMax {
						xMax = xs[i]
					}
					if a := abs64(xs[i]); a > ax {
						ax = a
					}
				}
			} else {
				if firstY {
					yMin, yMax = ys[i], ys[i]
					ay = abs64(ys[i])
					firstY = false
				} else {
					if ys[i] < yMin {
						yMin = ys[i]
					}
					if ys[i] > yMax {
						yMax = ys[i]
					}
					if a := abs64(ys[i]); a > ay {
						ay = a
					}
				}
			}
		}
		var dx2, dy2 int64
		if !firstX {
			d := xMax - xMin
			dx2 = d * d
		}
		if !firstY {
			d := yMax - yMin
			dy2 = d * d
		}
		cross := ax*ax + ay*ay
		res := dx2
		if dy2 > res {
			res = dy2
		}
		if cross > res {
			res = cross
		}
		return res
	}

	// Strategy 1: all to X
	assignAllX := make([]bool, n)
	for i := range assignAllX {
		assignAllX[i] = true
	}
	best := compute(assignAllX)

	// Strategy 2: all to Y
	assignAllY := make([]bool, n)
	for i := range assignAllY {
		assignAllY[i] = false
	}
	if v := compute(assignAllY); v < best {
		best = v
	}

	// Strategy 3: assign to axis with smaller absolute coordinate
	assignHeur := make([]bool, n)
	for i := 0; i < n; i++ {
		if abs64(xs[i]) <= abs64(ys[i]) {
			assignHeur[i] = true
		} else {
			assignHeur[i] = false
		}
	}
	if v := compute(assignHeur); v < best {
		best = v
	}

	fmt.Println(best)
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}
