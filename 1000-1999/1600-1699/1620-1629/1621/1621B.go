package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	const INF int64 = 1 << 60
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var lmin, rmax int64
		lmin = INF
		rmax = -INF
		costL := INF
		costR := INF
		best := INF
		for i := 0; i < n; i++ {
			var l, r, c int64
			fmt.Fscan(in, &l, &r, &c)
			if l < lmin {
				lmin = l
				costL = c
				best = INF
			} else if l == lmin && c < costL {
				costL = c
			}
			if r > rmax {
				rmax = r
				costR = c
				best = INF
			} else if r == rmax && c < costR {
				costR = c
			}
			if l == lmin && r == rmax && c < best {
				best = c
			}
			ans := costL + costR
			if best < ans {
				ans = best
			}
			if i > 0 {
				out.WriteByte('\n')
			}
			fmt.Fprint(out, ans)
		}
		out.WriteByte('\n')
	}
}
