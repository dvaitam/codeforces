package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	const INF = int(1e9)
	xmin, ymin := INF, INF
	xmax, ymax := -INF, -INF
	a, b := -INF, -INF

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		ctb, ctn := 0, 0
		for _, c := range s {
			if c == 'B' {
				ctb++
			} else {
				ctn++
			}
		}
		xmin = min(xmin, ctb)
		ymin = min(ymin, ctn)
		xmax = max(xmax, ctb)
		ymax = max(ymax, ctn)
		a = max(a, -ctb+ctn)
		b = max(b, ctb-ctn)
	}

	best, ansb, ansn := INF, 0, 0

	var check func(int) bool
	check = func(t int) bool {
		for xpr := 0; xpr <= xmax; xpr++ {
			if abs(xpr-xmin) > t || abs(xpr-xmax) > t {
				continue
			}
			ylow := ymax - t
			yhi := ymin + t
			ylow = max(ylow, xpr+a-t)
			yhi = min(yhi, t-b+xpr)
			if ylow <= yhi && (xpr > 0 || yhi > 0) {
				if t < best {
					best = t
					ansb = xpr
					ansn = yhi
				}
				return true
			}
		}
		return false
	}

	lo, hi := 0, int(1e6)
	for lo <= hi {
		mid := (lo + hi) / 2
		if check(mid) {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	fmt.Fprintln(writer, best)
	res := strings.Repeat("B", ansb) + strings.Repeat("N", ansn)
	fmt.Fprintln(writer, res)
}
