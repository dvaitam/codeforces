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

func get(dp [][]int64, l, r int) int64 {
	if l > r {
		return 0
	}
	return dp[l][r]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		dp := make([][]int64, n)
		for i := range dp {
			dp[i] = make([]int64, n)
		}

		for length := 1; length <= n; length++ {
			for l := 0; l+length-1 < n; l++ {
				r := l + length - 1
				best := int64(0)
				if l+1 <= r {
					best = dp[l+1][r]
				}
				for j := l + 2; j <= r; j++ {
					bestIJ := int64(-1) << 60
					for i := l + 1; i < j; i++ {
						val := get(dp, l+1, i-1) + get(dp, i+1, j-1) + int64(a[l])*int64(a[i])*int64(a[j])
						if val > bestIJ {
							bestIJ = val
						}
					}
					if bestIJ > int64(-1)<<50 {
						total := bestIJ + get(dp, j+1, r)
						if total > best {
							best = total
						}
					}
				}
				dp[l][r] = best
			}
		}

		ans := get(dp, 1, n-1)
		for i := 1; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				val := int64(a[0])*int64(a[i])*int64(a[j]) + get(dp, 1, i-1) + get(dp, i+1, j-1) + get(dp, j+1, n-1)
				if val > ans {
					ans = val
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
