package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

type constraint struct {
	pos int
	x   int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		cons := make([]constraint, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &cons[i].pos, &cons[i].x)
			cons[i].pos-- // zero-based
		}

		// val[s][pos] = number of valid rotations for block length s starting at pos
		val := make([][]int, n+1)
		for s := 1; s <= n; s++ {
			val[s] = make([]int, n-s+1)

			add := make([][]int, n) // residues for constraints at each index
			for _, c := range cons {
				if c.x > s {
					continue // never equals, constraint irrelevant
				}
				r := (c.x - (c.pos + 1)) % s
				if r < 0 {
					r += s
				}
				add[c.pos] = append(add[c.pos], r)
			}

			freq := make([]int, s)
			distinct := 0
			// initialize window [0, s-1]
			limit := s
			if limit > n {
				limit = n
			}
			for idx := 0; idx < limit; idx++ {
				for _, r := range add[idx] {
					if freq[r] == 0 {
						distinct++
					}
					freq[r]++
				}
			}

			for pos := 0; pos <= n-s; pos++ {
				val[s][pos] = s - distinct
				if pos == n-s {
					break
				}
				// slide window: remove pos, add pos+s
				for _, r := range add[pos] {
					freq[r]--
					if freq[r] == 0 {
						distinct--
					}
				}
				newIdx := pos + s
				if newIdx < n {
					for _, r := range add[newIdx] {
						if freq[r] == 0 {
							distinct++
						}
						freq[r]++
					}
				}
			}
		}

		dp := make([]int, n+1)
		dp[0] = 1
		for i := 0; i < n; i++ {
			if dp[i] == 0 {
				continue
			}
			for s := 1; i+s <= n; s++ {
				v := val[s][i]
				if v == 0 {
					continue
				}
				dp[i+s] = (dp[i+s] + int((int64(dp[i])*int64(v))%mod)) % mod
			}
		}

		fmt.Fprintln(out, dp[n]%mod)
	}
}
