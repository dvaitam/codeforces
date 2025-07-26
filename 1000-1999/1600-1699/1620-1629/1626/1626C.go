package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		k := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &k[i])
		}
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}

		// process intervals from the last monster backwards
		curL := k[n-1] - h[n-1] + 1
		curR := k[n-1]
		var ans int64
		for i := n - 2; i >= 0; i-- {
			start := k[i] - h[i] + 1
			end := k[i]
			if end >= curL { // overlap with current interval
				if start < curL {
					curL = start
				}
			} else {
				length := curR - curL + 1
				ans += length * (length + 1) / 2
				curL = start
				curR = end
			}
		}
		length := curR - curL + 1
		ans += length * (length + 1) / 2
		fmt.Fprintln(out, ans)
	}
}
