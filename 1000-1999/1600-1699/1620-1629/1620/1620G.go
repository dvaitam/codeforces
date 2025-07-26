package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	cnt := make([][26]int16, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for _, ch := range s {
			cnt[i][ch-'a']++
		}
	}
	m := 1 << n
	g := make([]int64, m)
	for i := 1; i < m; i++ {
		g[i] = 1
	}
	dp := make([]int16, m)
	for l := 0; l < 26; l++ {
		dp[0] = 0
		for mask := 1; mask < m; mask++ {
			lb := mask & -mask
			idx := bits.TrailingZeros(uint(lb))
			if mask == lb {
				dp[mask] = cnt[idx][l]
			} else {
				prev := mask ^ lb
				c := cnt[idx][l]
				if dp[prev] < c {
					dp[mask] = dp[prev]
				} else {
					dp[mask] = c
				}
			}
			g[mask] = g[mask] * int64(dp[mask]+1) % MOD
		}
	}
	g[0] = 0
	arr := make([]int64, m)
	for mask := 0; mask < m; mask++ {
		if bits.OnesCount(uint(mask))%2 == 1 {
			arr[mask] = (MOD - g[mask]) % MOD
		} else {
			arr[mask] = g[mask]
		}
	}
	for i := 0; i < n; i++ {
		for mask := 0; mask < m; mask++ {
			if mask&(1<<i) != 0 {
				arr[mask] = (arr[mask] + arr[mask^(1<<i)]) % MOD
			}
		}
	}
	idxSum := make([]int64, m)
	pop := make([]int, m)
	for mask := 1; mask < m; mask++ {
		lb := mask & -mask
		idx := bits.TrailingZeros(uint(lb))
		prev := mask ^ lb
		idxSum[mask] = idxSum[prev] + int64(idx+1)
		pop[mask] = pop[prev] + 1
	}
	var ans int64
	for mask := 1; mask < m; mask++ {
		f := arr[mask]
		if f != 0 {
			f = (MOD - f) % MOD
		}
		val := f * int64(pop[mask]) % MOD
		val = val * idxSum[mask]
		ans ^= val
	}
	fmt.Println(ans)
}
