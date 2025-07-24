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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)

		diff := make([]int, 0)
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				diff = append(diff, i)
			}
		}
		k := len(diff)
		if k%2 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		if k == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		if x >= y {
			if k == 2 && diff[1] == diff[0]+1 {
				ans := x
				if 2*y < ans {
					ans = 2 * y
				}
				fmt.Fprintln(out, ans)
			} else {
				fmt.Fprintln(out, int64(k/2)*y)
			}
		} else {
			// not expected in this version
			// fallback: dynamic programming
			dp := make([]int64, k+1)
			for i := range dp {
				dp[i] = 1 << 60
			}
			dp[0] = 0
			for i := 2; i <= k; i += 2 {
				// pair i-2 with i-1
				if diff[i-1]-diff[i-2] == 1 {
					if dp[i-2]+x < dp[i] {
						dp[i] = dp[i-2] + x
					}
				}
				// pair cross with other j
				for j := 0; j <= i-2; j += 2 {
					cost := dp[j] + int64((i-j)/2)*y
					if cost < dp[i] {
						dp[i] = cost
					}
				}
			}
			fmt.Fprintln(out, dp[k])
		}
	}
}
