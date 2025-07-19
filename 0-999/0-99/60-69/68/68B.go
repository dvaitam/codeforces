package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int
	fmt.Fscan(in, &n, &k)
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		var xi float64
		fmt.Fscan(in, &xi)
		x[i] = xi
	}
	low, high := 0.0, 1000.0
	for it := 0; it < 100; it++ {
		mid := (low + high) / 2
		if check(mid, x, float64(k)) {
			low = mid
		} else {
			high = mid
		}
	}
	res := (low + high) / 2
	fmt.Printf("%.8f", res)
}

// check returns true if total profit b at price a is positive
func check(a float64, x []float64, k float64) bool {
	b := 0.0
	for _, xi := range x {
		if xi > a {
			b += (xi - a) / 100 * (100 - k)
		} else {
			b -= a - xi
		}
	}
	return b > 0
}
