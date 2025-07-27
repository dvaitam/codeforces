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
		fmt.Fscan(in, &n)
		tags := make([]int, n+1)
		scores := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &tags[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &scores[i])
		}
		dp := make([]int64, n+1)
		var res int64
		for m := 1; m <= n; m++ {
			for j := m - 1; j >= 1; j-- {
				if tags[m] == tags[j] {
					continue
				}
				diffScore := abs64(scores[m] - scores[j])
				oldJ := dp[j]
				if dp[m]+diffScore > dp[j] {
					dp[j] = dp[m] + diffScore
				}
				if oldJ+diffScore > dp[m] {
					dp[m] = oldJ + diffScore
				}
			}
			if dp[m] > res {
				res = dp[m]
			}
		}
		for i := 1; i <= n; i++ {
			if dp[i] > res {
				res = dp[i]
			}
		}
		fmt.Fprintln(out, res)
	}
}
func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
