package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isqrt(x int64) int64 {
	if x < 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		freq := make(map[int64]int64, n*2)
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			freq[v]++
		}
		var q int
		fmt.Fscan(in, &q)
		for i := 0; i < q; i++ {
			var x, y int64
			fmt.Fscan(in, &x, &y)
			delta := x*x - 4*y
			var ans int64
			if delta >= 0 {
				s := isqrt(delta)
				if s*s == delta {
					if (x+s)%2 == 0 && (x-s)%2 == 0 {
						r1 := (x + s) / 2
						r2 := (x - s) / 2
						if r1 == r2 {
							c := freq[r1]
							if c >= 2 {
								ans = c * (c - 1) / 2
							}
						} else {
							ans = freq[r1] * freq[r2]
						}
					}
				}
			}
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
