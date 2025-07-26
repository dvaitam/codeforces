package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k, a, b int
		fmt.Fscan(in, &n, &k, &a, &b)
		xs := make([]int64, n+1)
		ys := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &xs[i], &ys[i])
		}
		direct := abs(xs[a]-xs[b]) + abs(ys[a]-ys[b])
		const inf int64 = 1 << 62
		dA, dB := inf, inf
		for i := 1; i <= k; i++ {
			da := abs(xs[a]-xs[i]) + abs(ys[a]-ys[i])
			if da < dA {
				dA = da
			}
			db := abs(xs[b]-xs[i]) + abs(ys[b]-ys[i])
			if db < dB {
				dB = db
			}
		}
		if dA+dB < direct {
			fmt.Fprintln(out, dA+dB)
		} else {
			fmt.Fprintln(out, direct)
		}
	}
}
