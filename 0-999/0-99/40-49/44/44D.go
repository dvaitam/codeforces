package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// Node represents a point in 3D space with its distance from the starting point
type Node struct {
	x, y, z, d float64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pts := make([]Node, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y, &pts[i].z)
	}
	// Compute initial distances from pts[0]
	for i := 0; i < n; i++ {
		dx := pts[i].x - pts[0].x
		dy := pts[i].y - pts[0].y
		dz := pts[i].z - pts[0].z
		pts[i].d = math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
	// Sort by distance
	sort.Slice(pts, func(i, j int) bool {
		return pts[i].d < pts[j].d
	})
	// Find minimal meeting time
	ans := -1.0
	for i := 1; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := pts[i].x - pts[j].x
			dy := pts[i].y - pts[j].y
			dz := pts[i].z - pts[j].z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
			// adjust by difference in initial times
			dist -= pts[j].d - pts[i].d
			when := pts[j].d + dist/2
			if ans < 0 || when < ans {
				ans = when
			}
		}
	}
	// Output with high precision
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintf(out, "%.19f", ans)
}
