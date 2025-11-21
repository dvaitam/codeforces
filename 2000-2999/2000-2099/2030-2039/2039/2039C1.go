package main

import (
	"bufio"
	"fmt"
	"os"
)

func divisors(n int) []int {
	res := []int{}
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			res = append(res, d)
			if d != n/d {
				res = append(res, n/d)
			}
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		var m int64
		fmt.Fscan(in, &x, &m)

		seen := make(map[int64]struct{})

		// Case 1: xor is a divisor of x
		for _, d := range divisors(x) {
			y := int64(x ^ d)
			if y <= 0 || y == int64(x) || y > m {
				continue
			}
			seen[y] = struct{}{}
		}

		// Case 2: xor is a divisor of y (d divides y)
		for d := 1; d <= x; d++ {
			yInt := x ^ d
			if yInt <= 0 || yInt == x {
				continue
			}
			y := int64(yInt)
			if y > m {
				continue
			}
			if yInt%d == 0 {
				seen[y] = struct{}{}
			}
		}

		fmt.Fprintln(out, len(seen))
	}
}
