package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int = 1e9

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
		var s, a string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &a)

		dp := make([]int, 8)
		for i := 0; i < 8; i++ {
			dp[i] = inf
		}
		dp[0] = 0

		for i := 0; i < n; i++ {
			ndp := make([]int, 8)
			for j := 0; j < 8; j++ {
				ndp[j] = inf
			}
			for mask := 0; mask < 8; mask++ {
				if dp[mask] == inf {
					continue
				}
				for bit := 0; bit < 2; bit++ {
					cost := dp[mask]
					if (s[i] == '(' && bit == 1) || (s[i] == ')' && bit == 0) {
						cost++
					}
					if i >= 3 && a[i-3] == '1' {
						b0 := (mask >> 2) & 1
						b1 := (mask >> 1) & 1
						b2 := mask & 1
						if !((b0 == 0 && b1 == 1 && b2 == 0 && bit == 1) ||
							(b0 == 0 && b1 == 0 && b2 == 1 && bit == 1)) {
							continue
						}
					}
					nm := ((mask << 1) & 7) | bit
					if cost < ndp[nm] {
						ndp[nm] = cost
					}
				}
			}
			dp = ndp
		}

		ans := inf
		for _, v := range dp {
			if v < ans {
				ans = v
			}
		}
		if ans == inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
