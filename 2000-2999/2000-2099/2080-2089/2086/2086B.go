package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; t > 0; t-- {
		var n, k int
		var x int64
		fmt.Fscan(in, &n, &k, &x)

		a := make([]int64, n)
		var total int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
		}

		suffix := make([]int64, n)
		var running int64
		for i := n - 1; i >= 0; i-- {
			running += a[i]
			suffix[i] = running
		}

		var ans int64
		k64 := int64(k)
		for i := 0; i < n; i++ {
			need := x - suffix[i]
			if need <= 0 {
				ans += k64
				continue
			}

			// Positive numbers guarantee total > 0.
			m := (need + total - 1) / total
			if m >= k64 {
				// No copy index satisfies the requirement.
				continue
			}
			ans += k64 - m
		}

		fmt.Fprintln(out, ans)
	}
}
