package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	dp := make([]int64, k)
	sum := make([]int64, k)
	dp[0] = 1

	prefixDp := make([]int64, k+1)
	prefixSum := make([]int64, k+1)
	nextDp := make([]int64, k)
	nextSum := make([]int64, k)

	for step := 0; step < n; step++ {
		// compute prefix sums for repetition transitions
		prefixDp[k] = 0
		prefixSum[k] = 0
		for i := k - 1; i >= 0; i-- {
			prefixDp[i] = prefixDp[i+1] + dp[i]
			if prefixDp[i] >= mod {
				prefixDp[i] %= mod
			}
			prefixSum[i] = prefixSum[i+1] + sum[i]
			if prefixSum[i] >= mod {
				prefixSum[i] %= mod
			}
		}

		for i := 0; i < k; i++ {
			nextDp[i] = 0
			nextSum[i] = 0
		}

		// contributions from repeating digits
		for r := 1; r < k; r++ {
			nextDp[r] = prefixDp[r]
			nextSum[r] = prefixSum[r]
		}

		// contributions from new distinct digits
		for m := 0; m < k; m++ {
			if dp[m] == 0 && sum[m] == 0 {
				continue
			}
			newChoices := int64(k - m)
			if m+1 == k {
				nextDp[0] = (nextDp[0] + dp[m]*newChoices) % mod
				inc := (sum[m] + dp[m]) % mod
				nextSum[0] = (nextSum[0] + inc*newChoices) % mod
			} else {
				nextDp[m+1] = (nextDp[m+1] + dp[m]*newChoices) % mod
				nextSum[m+1] = (nextSum[m+1] + sum[m]*newChoices) % mod
			}
		}

		dp, nextDp = nextDp, dp
		sum, nextSum = nextSum, sum
	}

	var ans int64
	for i := 0; i < k; i++ {
		ans = (ans + sum[i]) % mod
	}
	fmt.Fprintln(out, ans)
}
