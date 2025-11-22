package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type point struct {
	x, y int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	pts := make([]point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if k == 1 {
		fmt.Fprintln(out, "0.000000000000000")
		return
	}

	if n == 1 {
		fmt.Fprintln(out, "0.000000000000000")
		return
	}

	// All points must be collinear (no three points can be on a common circle of finite radius).
	// Project onto the direction of the first segment to reduce to 1D.
	dx := pts[1].x - pts[0].x
	dy := pts[1].y - pts[0].y
	norm := math.Hypot(float64(dx), float64(dy))

	proj := make([]int64, n)
	for i := 0; i < n; i++ {
		proj[i] = (pts[i].x-pts[0].x)*dx + (pts[i].y-pts[0].y)*dy
	}
	sort.Slice(proj, func(i, j int) bool { return proj[i] < proj[j] })

	minDiff := proj[k-1] - proj[0]
	for i := k; i < n; i++ {
		if diff := proj[i] - proj[i-k+1]; diff < minDiff {
			minDiff = diff
		}
	}

	radius := float64(minDiff) / (2.0 * norm)
	area := math.Pi * radius * radius
	fmt.Fprintf(out, "%.15f\n", area)
}
