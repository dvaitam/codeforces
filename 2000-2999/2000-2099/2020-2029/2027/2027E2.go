package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

// calcCounts returns counts[c] = # of x in [0, b] with popcount(a & x) == c (mod mod).
func calcCounts(a, b int64) [31]int64 {
	var tight, loose [31]int64
	tight[0] = 1
	for bit := 29; bit >= 0; bit-- {
		aBit := (a >> bit) & 1
		bBit := (b >> bit) & 1
		var ntight, nloose [31]int64

		// From tight states.
		for c := 0; c <= 30; c++ {
			val := tight[c]
			if val == 0 {
				continue
			}
			if bBit == 0 {
				// Must place 0 to stay <= b.
				ntight[c] = (ntight[c] + val) % mod
			} else {
				// Place 0 -> becomes loose.
				nloose[c] = (nloose[c] + val) % mod
				// Place 1 -> stays tight.
				c2 := c + int(aBit)
				ntight[c2] = (ntight[c2] + val) % mod
			}
		}

		// From loose states (already < b, can place 0/1 freely).
		for c := 0; c <= 30; c++ {
			val := loose[c]
			if val == 0 {
				continue
			}
			// Place 0.
			nloose[c] = (nloose[c] + val) % mod
			// Place 1.
			c2 := c + int(aBit)
			nloose[c2] = (nloose[c2] + val) % mod
		}

		tight, loose = ntight, nloose
	}

	var res [31]int64
	for i := 0; i <= 30; i++ {
		res[i] = (tight[i] + loose[i]) % mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		dp := make([]int64, 31)
		dp[0] = 1

		for i := 0; i < n; i++ {
			cnt := calcCounts(a[i], b[i])
			// Exclude x = 0 (not allowed).
			cnt[0] = (cnt[0] + mod - 1) % mod

			ndp := make([]int64, 31)
			for x := 0; x <= 30; x++ {
				if dp[x] == 0 {
					continue
				}
				for c := 0; c <= 30; c++ {
					if cnt[c] == 0 {
						continue
					}
					nx := x ^ c
					ndp[nx] = (ndp[nx] + dp[x]*cnt[c]) % mod
				}
			}
			dp = ndp
		}

		fmt.Fprintln(out, dp[0]%mod)
	}
}
