package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func ceilSqrt(n int64) int {
	if n <= 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(n)))
	for r*r < n {
		r++
	}
	for (r-1)*(r-1) >= n {
		r--
	}
	return int(r)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var X, Y, K int
		fmt.Fscan(in, &X, &Y, &K)

		m := X
		if Y < m {
			m = Y
		}

		target := int64(K) * int64(K)
		dx, dy := -1, -1

		for cand := 0; cand <= m; cand++ {
			need := target - int64(cand)*int64(cand)
			curDy := 0
			if need > 0 {
				curDy = ceilSqrt(need)
			}
			if curDy > m {
				continue
			}
			if cand == 0 && curDy == 0 {
				continue
			}
			dx = cand
			dy = curDy
			break
		}

		if dx == -1 {
			// As per problem statement this shouldn't happen, but keep fallback.
			fmt.Fprintln(out, "0 0 0 0")
			fmt.Fprintln(out, "0 0 0 0")
			continue
		}

		// First segment: (0,0) to (dx, dy)
		fmt.Fprintf(out, "0 0 %d %d\n", dx, dy)
		// Second segment: (dy,0) to (0,dx), perpendicular to the first.
		fmt.Fprintf(out, "%d 0 0 %d\n", dy, dx)
	}
}
