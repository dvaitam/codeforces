package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var r0, x, d, n int
		fmt.Fscan(in, &r0, &x, &d, &n)
		var rounds string
		fmt.Fscan(in, &rounds)

		low := int64(r0)
		high := int64(r0)
		delta := int64(d)
		threshold := int64(x)
		rated := 0

		for i := 0; i < n; i++ {
			if rounds[i] == '1' {
				rated++
				low = max(0, low-delta)
				high = high + delta
				continue
			}

			if low < threshold {
				rated++
				upperStart := min(high, threshold-1)
				low = max(0, low-delta)
				high = upperStart + delta
			}
		}

		fmt.Fprintln(out, rated)
	}
}
