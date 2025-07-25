package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		type interval struct{ l, r int }
		intervals := make([]interval, 0, 2*n)
		for i := 1; i <= n; i++ {
			l1 := i - a[i] + 1
			if l1 < 1 {
				l1 = 1
			}
			intervals = append(intervals, interval{l1, i})
			r2 := i + a[i] - 1
			if r2 > n {
				r2 = n
			}
			intervals = append(intervals, interval{i, r2})
		}
		const INF = int(1e9)
		dp := make([]int, n+2)
		for i := 0; i <= n; i++ {
			dp[i] = INF
		}
		dp[n+1] = 0
		for i := n; i >= 1; i-- {
			best := INF
			for _, seg := range intervals {
				if seg.l <= i && seg.r >= i {
					if val := dp[seg.r+1] + 1; val < best {
						best = val
					}
				}
			}
			dp[i] = best
		}
		fmt.Fprintln(writer, dp[1])
	}
}
