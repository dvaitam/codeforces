package main

import (
	"bufio"
	"fmt"
	"os"
)

// Point represents a coordinate pair
type Point struct {
	x, y int64
}

// dist returns squared distance between two points
func dist(a, b Point) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return dx*dx + dy*dy
}

// permute generates all permutations of a and updates best sequence for max path
func permute(a []Point, l int, bestSum *int64, bestSeq *[]Point) {
	if l == len(a)-1 {
		// compute path sum of consecutive distances
		sum := dist(a[0], a[1]) + dist(a[1], a[2]) + dist(a[2], a[3])
		if sum > *bestSum {
			*bestSum = sum
			*bestSeq = make([]Point, 4)
			copy(*bestSeq, a)
		}
	} else {
		for i := l; i < len(a); i++ {
			a[l], a[i] = a[i], a[l]
			permute(a, l+1, bestSum, bestSeq)
			a[l], a[i] = a[i], a[l]
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	// special cases
	if n == 0 {
		fmt.Fprintf(writer, "%d %d\n", 0, 1)
		fmt.Fprintf(writer, "%d %d\n", 0, m)
		fmt.Fprintf(writer, "%d %d\n", 0, 0)
		fmt.Fprintf(writer, "%d %d\n", 0, m-1)
		return
	}
	if m == 0 {
		fmt.Fprintf(writer, "%d %d\n", 1, 0)
		fmt.Fprintf(writer, "%d %d\n", n, 0)
		fmt.Fprintf(writer, "%d %d\n", 0, 0)
		fmt.Fprintf(writer, "%d %d\n", n-1, 0)
		return
	}
	// candidate points
	pts := []Point{
		{0, 0}, {n, 0}, {0, m}, {n, m},
		{1, 0}, {0, 1}, {n - 1, 0}, {0, m - 1},
		{n, 1}, {1, m}, {n - 1, m}, {n, m - 1},
	}
	// dedupe
	uniq := make(map[Point]bool)
	candidates := make([]Point, 0, len(pts))
	for _, p := range pts {
		if !uniq[p] {
			uniq[p] = true
			candidates = append(candidates, p)
		}
	}
	// search best sequence
	var bestSum int64 = -1
	var bestSeq []Point
	L := len(candidates)
	for i := 0; i < L; i++ {
		for j := i + 1; j < L; j++ {
			for k := j + 1; k < L; k++ {
				for l := k + 1; l < L; l++ {
					subset := []Point{candidates[i], candidates[j], candidates[k], candidates[l]}
					permute(subset, 0, &bestSum, &bestSeq)
				}
			}
		}
	}
	// output
	for _, p := range bestSeq {
		fmt.Fprintf(writer, "%d %d\n", p.x, p.y)
	}
}
