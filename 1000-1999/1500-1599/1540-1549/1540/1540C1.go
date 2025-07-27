package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}
	b := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &b[i])
	}
	var q int
	fmt.Fscan(in, &q)
	xs := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &xs[i])
	}

	// Precompute prefix sums for b and weights w
	prefixB := make([]int64, n)
	for i := 1; i < n; i++ {
		prefixB[i] = prefixB[i-1] + int64(b[i-1])
	}
	w := make([]int64, n)
	for i := 1; i < n; i++ {
		w[i] = w[i-1] + prefixB[i]
	}

	maxSum := 0
	for _, v := range c {
		maxSum += v
	}

	for _, x := range xs {
		dp := make([]int64, maxSum+1)
		dp[0] = 1
		limits := make([]int64, n)
		for i := 0; i < n; i++ {
			limits[i] = w[i] + int64(i+1)*int64(x)
		}
		for i := 0; i < n; i++ {
			ndp := make([]int64, maxSum+1)
			pref := make([]int64, maxSum+1)
			pref[0] = dp[0]
			for s := 1; s <= maxSum; s++ {
				pref[s] = (pref[s-1] + dp[s]) % MOD
			}
			low := int(limits[i])
			if low > maxSum {
				dp = ndp
				break
			}
			if low < 0 {
				low = 0
			}
			for s := low; s <= maxSum; s++ {
				l := s - c[i]
				if l < 0 {
					l = 0
				}
				val := pref[s]
				if l > 0 {
					val = (val - pref[l-1]) % MOD
				}
				if val < 0 {
					val += MOD
				}
				ndp[s] = val
			}
			dp = ndp
		}
		ans := int64(0)
		for _, v := range dp {
			ans = (ans + v) % MOD
		}
		fmt.Println(ans)
	}
}
