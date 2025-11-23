package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	type point struct {
		x, y int64
	}

	raw := make([]point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &raw[i].x, &raw[i].y)
	}

	// Sort points by x coordinate, then by y coordinate.
	sort.Slice(raw, func(i, j int) bool {
		if raw[i].x != raw[j].x {
			return raw[i].x < raw[j].x
		}
		return raw[i].y < raw[j].y
	})

	var hull []point
	for i := 0; i < n; i++ {
		// If the current point has the same x as the next one, it is not the highest point
		// for this x coordinate (due to sorting). We only care about the highest point.
		if i+1 < n && raw[i].x == raw[i+1].x {
			continue
		}

		// Transform coordinates: Y = y - x^2, X = x
		// We are looking for the upper convex hull of these transformed points.
		p := point{raw[i].x, raw[i].y - raw[i].x*raw[i].x}

		for len(hull) >= 2 {
			a := hull[len(hull)-2]
			b := hull[len(hull)-1]
			
			// Check for convexity. We want the upper hull.
			// Slopes must be strictly decreasing.
			// If slope(a, b) <= slope(b, p), then b is below or on the segment a-p,
			// so we remove b.
			// Inequality: (b.y - a.y) / (b.x - a.x) <= (p.y - b.y) / (p.x - b.x)
			// Cross-multiplication (denominators are positive because sorted by x):
			// (b.y - a.y) * (p.x - b.x) <= (p.y - b.y) * (b.x - a.x)
			
			lhs := (b.y - a.y) * (p.x - b.x)
			rhs := (p.y - b.y) * (b.x - a.x)

			if lhs <= rhs {
				hull = hull[:len(hull)-1]
			} else {
				break
			}
		}
		hull = append(hull, p)
	}

	ans := 0
	if len(hull) > 1 {
		ans = len(hull) - 1
	}
	fmt.Fprintln(out, ans)
}