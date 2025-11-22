package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod int64 = 1_000_000_007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func removeBit(mask int, idx int) int {
	low := mask & ((1 << idx) - 1)
	high := mask >> (idx + 1)
	return low | (high << idx)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(in, &n, &m)

		var k int
		fmt.Fscan(in, &k)
		goodIdx := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &goodIdx[i])
			goodIdx[i]-- // zero-based
		}

		// precompute factorials up to n for combinations modulo mod
		fact := make([]int64, n+1)
		invFact := make([]int64, n+1)
		fact[0] = 1
		for i := 1; i <= n; i++ {
			fact[i] = fact[i-1] * int64(i) % mod
		}
		invFact[n] = modPow(fact[n], mod-2)
		for i := n; i > 0; i-- {
			invFact[i-1] = invFact[i] * int64(i) % mod
		}
		combMod := func(nn, rr int) int64 {
			if rr < 0 || rr > nn {
				return 0
			}
			return fact[nn] * invFact[rr] % mod * invFact[nn-rr] % mod
		}

		// dp[s][mask] -> bool outcome with s piles represented by lower s bits of mask
		dp := make([][]bool, n+1)
		dp[1] = make([]bool, 2)
		dp[1][0] = false
		dp[1][1] = true

		for s := 2; s <= n; s++ {
			size := 1 << s
			dp[s] = make([]bool, size)
			turnAlice := ((n - s) & 1) == 0 // true if Alice (maximizer) plays
			for mask := 0; mask < size; mask++ {
				var res bool
				if turnAlice {
					res = false
					for _, g := range goodIdx {
						if g >= s {
							break
						}
						child := removeBit(mask, g)
						if dp[s-1][child] {
							res = true
							break
						}
					}
				} else {
					res = true
					for _, g := range goodIdx {
						if g >= s {
							break
						}
						child := removeBit(mask, g)
						if !dp[s-1][child] {
							res = false
							break
						}
					}
				}
				dp[s][mask] = res
			}
		}

		// count masks of length n with popcount t that result in black
		cnt := make([]int64, n+1)
		size := 1 << n
		for mask := 0; mask < size; mask++ {
			if dp[n][mask] {
				pop := bits.OnesCount(uint(mask))
				cnt[pop]++
			}
		}

		expectedRank := int64(0)
		for t := 1; t <= n; t++ {
			ways := combMod(n, t)
			if ways == 0 {
				continue
			}
			expectedRank = (expectedRank + cnt[t]%mod*modPow(ways, mod-2)) % mod
		}

		// E[value] = (m+1)/(n+1) * expectedRank
		ev := (expectedRank * ((m + 1) % mod)) % mod
		ev = (ev * modPow(int64(n+1), mod-2)) % mod

		totalWays := modPow(m%mod, int64(n))
		ans := ev * totalWays % mod
		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(out, ans)
	}
}
