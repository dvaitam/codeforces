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
	x := make([]float64, n)
	v := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &v[i])
	}

	ok := func(t float64) bool {
		left := -1e18
		right := 1e18
		for i := 0; i < n; i++ {
			l := x[i] - v[i]*t
			r := x[i] + v[i]*t
			if l > left {
				left = l
			}
			if r < right {
				right = r
			}
		}
		return left <= right
	}

	lo, hi := 0.0, 1e9
	for iter := 0; iter < 100; iter++ {
		mid := (lo + hi) / 2
		if ok(mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	fmt.Printf("%.10f\n", hi)
}
