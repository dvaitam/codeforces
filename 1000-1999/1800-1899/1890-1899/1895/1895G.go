package main

import (
	"bufio"
	"fmt"
	"math"
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
		var s string
		fmt.Fscan(reader, &s)
		r := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &r[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		base := int64(0)
		for i := 0; i < n; i++ {
			base += b[i]
		}

		onesTotal := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				onesTotal++
			}
		}

		inf := int64(math.MinInt64 / 4)
		dp := make([]int64, onesTotal+1)
		for i := range dp {
			dp[i] = inf
		}
		dp[0] = 0
		onesSeen := 0

		for idx := 0; idx < n; idx++ {
			diff := r[idx] - b[idx]
			if s[idx] == '1' {
				for k := onesSeen; k >= 0; k-- {
					if dp[k] == inf {
						continue
					}
					cand := dp[k] + diff
					if cand > dp[k+1] {
						dp[k+1] = cand
					}
				}
				onesSeen++
			} else {
				for k := 0; k <= onesSeen; k++ {
					if dp[k] == inf {
						continue
					}
					cand := dp[k] + diff - int64(k)
					if cand > dp[k] {
						dp[k] = cand
					}
				}
			}
		}
		ans := int64(math.MinInt64)
		for k := 0; k <= onesSeen; k++ {
			if dp[k] > ans {
				ans = dp[k]
			}
		}
		fmt.Fprintln(writer, base+ans)
	}
}
