package main

import (
	"bufio"
	"fmt"
	"os"
)

func powMod(a, b, mod int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, m int64
	if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
		return
	}

	total := powMod(k, n, m)
	var countNoLucky int64

	for S := int64(0); S < k; S++ {
		// Build list of allowed digits for this S
		allowed := make([]int64, 0, k)
		for d := int64(0); d < k; d++ {
			if (2*d)%k != S {
				allowed = append(allowed, d)
			}
		}
		dp := make([]int64, k)
		dp[0] = 1
		for i := int64(0); i < n; i++ {
			next := make([]int64, k)
			for s := int64(0); s < k; s++ {
				if dp[s] == 0 {
					continue
				}
				for _, d := range allowed {
					ns := (s + d) % k
					next[ns] += dp[s]
					if next[ns] >= m {
						next[ns] %= m
					}
				}
			}
			dp = next
		}
		countNoLucky += dp[S]
		if countNoLucky >= m {
			countNoLucky %= m
		}
	}

	ans := (total - countNoLucky) % m
	if ans < 0 {
		ans += m
	}
	fmt.Fprintln(out, ans)
}
