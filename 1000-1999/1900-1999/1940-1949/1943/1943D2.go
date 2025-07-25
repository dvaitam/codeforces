// Solution for 1943D2 - Hard version

package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n, k int, mod int64) int64 {
	dpBuf := make([]int64, (k+1)*(k+1))
	dp := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = dpBuf[i*(k+1) : (i+1)*(k+1)]
	}
	for a1 := 0; a1 <= k; a1++ {
		for a2 := a1; a2 <= k; a2++ {
			dp[a1][a2] = 1
		}
	}

	nextBuf := make([]int64, (k+1)*(k+1))
	next := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		next[i] = nextBuf[i*(k+1) : (i+1)*(k+1)]
	}

	arrBuf := make([]int64, (k+1)*(k+2))
	arr := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		arr[i] = arrBuf[i*(k+2) : (i+1)*(k+2)]
	}

	for pos := 3; pos <= n; pos++ {
		for i := range nextBuf {
			nextBuf[i] = 0
		}
		for i := range arrBuf {
			arrBuf[i] = 0
		}
		for prev1 := 0; prev1 <= k; prev1++ {
			for prev2 := 0; prev2 <= k; prev2++ {
				cnt := dp[prev1][prev2]
				if cnt == 0 {
					continue
				}
				L := prev2 - prev1
				if L < 0 {
					L = 0
				}
				arr[prev2][L] = (arr[prev2][L] + cnt) % mod
				arr[prev2][k+1] = (arr[prev2][k+1] - cnt) % mod
			}
		}
		for prev2 := 0; prev2 <= k; prev2++ {
			cur := int64(0)
			for val := 0; val <= k; val++ {
				cur = (cur + arr[prev2][val]) % mod
				if cur < 0 {
					cur += mod
				}
				next[prev2][val] = cur
			}
		}
		dp, next = next, dp
	}

	ans := int64(0)
	for prev1 := 0; prev1 <= k; prev1++ {
		for prev2 := 0; prev2 <= k && prev2 <= prev1; prev2++ {
			ans += dp[prev1][prev2]
		}
		ans %= mod
	}
	return ans % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		var p int64
		fmt.Fscan(in, &n, &k, &p)
		res := solve(n, k, p)
		fmt.Fprintln(out, res)
	}
}
