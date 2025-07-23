package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	F := make([][]int64, n)
	for i := 0; i < n; i++ {
		F[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &F[i][j])
		}
	}

	S := make([][]int64, n)
	for i := 0; i < n; i++ {
		S[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &S[i][j])
		}
	}

	half := n / 2
	total := 1 << n
	dp := make([][]int64, total)
	for i := range dp {
		dp[i] = make([]int64, half+1)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	dp[0][0] = 0

	for mask := 0; mask < total; mask++ {
		i := bits.OnesCount(uint(mask))
		if i >= n {
			continue
		}
		for f := 0; f <= half; f++ {
			cur := dp[mask][f]
			if cur < 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if mask&(1<<j) != 0 {
					continue
				}
				nmask := mask | (1 << j)
				dp[nmask][f] = max64(dp[nmask][f], cur+S[i][j])
				if f < half {
					dp[nmask][f+1] = max64(dp[nmask][f+1], cur+F[i][j])
				}
			}
		}
	}

	fmt.Fprintln(out, dp[total-1][half])
}
