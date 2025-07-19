package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var pInt int64
	if _, err := fmt.Fscan(reader, &n, &pInt); err != nil {
		return
	}
	a := make([]float64, n)
	b := make([]float64, n)
	var sumA float64
	for i := 0; i < n; i++ {
		var ai, bi int64
		fmt.Fscan(reader, &ai, &bi)
		a[i] = float64(ai)
		b[i] = float64(bi)
		sumA += a[i]
	}
	p := float64(pInt)
	if sumA <= p {
		fmt.Fprintln(writer, -1)
		return
	}
	l, r := 0.0, 1e18
	// binary search for maximum time
	for it := 0; it < 100; it++ {
		mid := (l + r) / 2
		if feasible(mid, a, b, p) {
			l = mid
		} else {
			r = mid
		}
	}
	fmt.Fprintf(writer, "%.9f\n", l)
}

// feasible checks if time mid is sufficient
func feasible(mid float64, a, b []float64, p float64) bool {
	var needSum float64
	for i := range a {
		cnt := mid * a[i]
		if cnt > b[i] {
			needSum += cnt - b[i]
		}
	}
	return needSum <= mid*p
}
