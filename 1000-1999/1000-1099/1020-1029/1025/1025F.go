package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Point represents a point in 2D space
type Point struct {
	x, y int64
}

// quad determines which half-plane/quadrant the point belongs to relative to the origin.
// It divides the plane into two parts [0, pi) and [pi, 2*pi).
// Returns 1 for upper half-plane (including positive x-axis), 2 for lower.
func quad(p Point) int {
	if p.y > 0 || (p.y == 0 && p.x > 0) {
		return 1
	}
	return 2
}

// cross computes the cross product of two points (vectors)
func cross(a, b Point) int64 {
	return a.x*b.y - a.y*b.x
}

func main() {
	// Use buffered I/O for faster execution
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	points := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &points[i].x, &points[i].y)
	}

	var ans int64
	// Buffer for storing relative points to avoid reallocation inside the loop
	// Capacity is 2*n to accommodate the duplicated array
	rel := make([]Point, 0, 2*n)

	// Iterate over each point to use as the pivot/origin
	for i := 0; i < n; i++ {
		rel = rel[:0]
		origin := points[i]

		// Collect all other points relative to origin
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			rel = append(rel, Point{points[j].x - origin.x, points[j].y - origin.y})
		}

		// Sort points angularly around the origin
		sort.Slice(rel, func(a, b int) bool {
			qa := quad(rel[a])
			qb := quad(rel[b])
			if qa != qb {
				return qa < qb
			}
			return cross(rel[a], rel[b]) > 0
		})

		m := len(rel) // m = n - 1
		// Duplicate the slice to handle angular wrap-around easily
		for j := 0; j < m; j++ {
			rel = append(rel, rel[j])
		}

		// Use two pointers (sliding window) to find points strictly to the left
		k := 0
		for j := 0; j < m; j++ {
			// Ensure k is at least the next point
			if k <= j {
				k = j + 1
			}
			// Expand window while points are strictly to the left (angle < 180 degrees)
			// cross > 0 ensures strictly to the left since no three points are collinear
			for k < j+m && cross(rel[j], rel[k]) > 0 {
				k++
			}

			// Number of points strictly to the left of the directed line i -> rel[j]
			cntL := int64(k - j - 1)
			
			// Number of points strictly to the right
			// Total other points is m = n - 1. 
			// One is used as the end of the line (rel[j]), so remaining is m - 1.
			cntR := int64(m - 1) - cntL

			// If we can form triangles on both sides
			if cntL >= 2 && cntR >= 2 {
				// Choose 2 points from left set and 2 from right set
				waysL := cntL * (cntL - 1) / 2
				waysR := cntR * (cntR - 1) / 2
				ans += waysL * waysR
			}
		}
	}

	// Each pair of disjoint triangles is separated by exactly two internal tangents.
	// Thus, we have counted each pair exactly twice.
	fmt.Fprintln(writer, ans/2)
}