package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	y, a, b    float64
	n          int
	l, r, pref []float64
)

// solve computes the total covered length up to x using prefix sums and binary search
func solve(x float64) float64 {
	low, high := 0, n+1
	for high-low > 1 {
		mid := (low + high) / 2
		if r[mid] >= x {
			high = mid
		} else {
			low = mid
		}
	}
	tr := high
	var res float64
	if x > l[tr] {
		res = x - l[tr]
	}
	res += pref[tr-1]
	return res
}

// getIntersection computes intersection of line through (x1,y1)-(x2,y2) with horizontal line y=0
func getIntersection(x1, y1, x2, y2 float64) float64 {
	return x1 - y1*(x1-x2)/(y1-y2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &y, &a, &b)
	fmt.Fscan(reader, &n)
	l = make([]float64, n+2)
	r = make([]float64, n+2)
	pref = make([]float64, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &l[i], &r[i])
		pref[i] = pref[i-1] + (r[i] - l[i])
	}
	// sentinel intervals
	l[0], r[0] = -2e9, -2e9
	l[n+1], r[n+1] = 2e9, 2e9
	pref[n+1] = pref[n]

	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		var px, py float64
		fmt.Fscan(reader, &px, &py)
		ax := getIntersection(px, py, a, y)
		bx := getIntersection(px, py, b, y)
		k := (py - y) / py
		res := (solve(bx) - solve(ax)) * k
		fmt.Fprintf(writer, "%.15f\n", res)
	}
}
