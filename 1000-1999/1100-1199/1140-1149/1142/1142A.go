package main

import (
	"bufio"
	"fmt"
	"os"
)

// gcd returns the greatest common divisor of x and y.
func gcd(x, y int64) int64 {
	for x != 0 {
		y %= x
		x, y = y, x
	}
	return y
}

// abs returns the absolute value of x.
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k, a, b int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	fmt.Fscan(reader, &a, &b)

	mn := n * k
	var mx int64

	xs := []int64{-a, a}
	ys := []int64{-b, b}
	for _, x := range xs {
		for _, y := range ys {
			d := abs(y - x)
			by := gcd(k, d)

			rk := k / by
			rn := n
			for g := gcd(rk, rn); g > 1; g = gcd(rk, rn) {
				rn /= g
			}

			// minimal stops
			valMin := (n / rn) * (k / by)
			if valMin < mn {
				mn = valMin
			}
			// maximal stops
			valMax := (n * k) / by
			if valMax > mx {
				mx = valMax
			}
		}
	}

	fmt.Printf("%d %d\n", mn, mx)
}
