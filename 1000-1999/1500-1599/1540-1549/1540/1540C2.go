package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

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

	prefixB := make([]int64, n)
	for i := 1; i < n; i++ {
		prefixB[i] = prefixB[i-1] + int64(b[i-1])
	}
	B := make([]int64, n+1)
	for i := 2; i <= n; i++ {
		B[i] = B[i-1] + prefixB[i-1]
	}
	prefixC := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefixC[i] = prefixC[i-1] + int64(c[i-1])
	}
	totalWays := int64(1)
	for i := 0; i < n; i++ {
		totalWays = (totalWays * int64(c[i]+1)) % mod
	}

	lowBound := int64(1 << 60)
	for i := 1; i <= n; i++ {
		val := (-B[i]) / int64(i)
		if val < lowBound {
			lowBound = val
		}
	}
	highBound := int64(1 << 60)
	for i := 1; i <= n; i++ {
		val := (prefixC[i] - B[i]) / int64(i)
		if val < highBound {
			highBound = val
		}
	}

	cache := make(map[int]int64)

	var solve func(int) int64
	solve = func(x int) int64 {
		if int64(x) <= lowBound {
			return totalWays
		}
		if int64(x) > highBound {
			return 0
		}
		if v, ok := cache[x]; ok {
			return v
		}
		maxSum := int(prefixC[n])
		dp := make([]int64, maxSum+1)
		dp[0] = 1
		for i := 1; i <= n; i++ {
			pref := make([]int64, maxSum+1)
			pref[0] = dp[0]
			for s := 1; s <= maxSum; s++ {
				pref[s] = (pref[s-1] + dp[s]) % mod
			}
			ndp := make([]int64, maxSum+1)
			for s := 0; s <= maxSum; s++ {
				l := s - c[i-1]
				if l <= 0 {
					ndp[s] = pref[s] % mod
				} else {
					val := pref[s] - pref[l-1]
					val %= mod
					if val < 0 {
						val += mod
					}
					ndp[s] = val
				}
			}
			thresh := int(B[i] + int64(i)*int64(x))
			if thresh > 0 {
				if thresh > maxSum {
					return 0
				}
				for s := 0; s < thresh; s++ {
					ndp[s] = 0
				}
			}
			dp = ndp
		}
		thresh := int(B[n] + int64(n)*int64(x))
		if thresh < 0 {
			thresh = 0
		} else if thresh > maxSum {
			return 0
		}
		ans := int64(0)
		for s := thresh; s <= maxSum; s++ {
			ans += dp[s]
			if ans >= mod {
				ans -= mod
			}
		}
		cache[x] = ans % mod
		return ans % mod
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, solve(xs[i]))
	}
}
