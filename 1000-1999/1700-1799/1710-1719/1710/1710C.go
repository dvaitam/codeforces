package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	nBits := len(s)
	// dp[state][mask] => number of ways for processed prefix
	// state: 3 bits tight flags for a,b,c (0 equal,1 less)
	// mask: bits of conditions satisfied so far
	dp := make([][64]int, 8)
	dp[0][0] = 1

	// precompute event mask for each combination of bits a,b,c
	eventMask := make([]int, 8)
	for i := 0; i < 8; i++ {
		a := (i >> 2) & 1
		b := (i >> 1) & 1
		c := i & 1
		m := 0
		if a != b {
			m |= 1 << 0 // diff_ab
		}
		if b != c {
			m |= 1 << 1 // diff_bc
		}
		if a != c {
			m |= 1 << 2 // diff_ac
		}
		if a != b && b != c {
			m |= 1 << 3 // ab & bc
		}
		if a != b && a != c {
			m |= 1 << 4 // ab & ac
		}
		if b != c && a != c {
			m |= 1 << 5 // bc & ac
		}
		eventMask[i] = m
	}

	for idx := 0; idx < nBits; idx++ {
		bit := int(s[idx] - '0')
		next := make([][64]int, 8)
		for st := 0; st < 8; st++ {
			ta := (st >> 2) & 1
			tb := (st >> 1) & 1
			tc := st & 1
			for mask := 0; mask < 64; mask++ {
				cur := dp[st][mask]
				if cur == 0 {
					continue
				}
				for combo := 0; combo < 8; combo++ {
					aBit := (combo >> 2) & 1
					bBit := (combo >> 1) & 1
					cBit := combo & 1
					// check against n bit
					if ta == 0 && aBit > bit {
						continue
					}
					if tb == 0 && bBit > bit {
						continue
					}
					if tc == 0 && cBit > bit {
						continue
					}
					nta := ta
					ntb := tb
					ntc := tc
					if ta == 0 && aBit < bit {
						nta = 1
					}
					if tb == 0 && bBit < bit {
						ntb = 1
					}
					if tc == 0 && cBit < bit {
						ntc = 1
					}
					nst := (nta << 2) | (ntb << 1) | ntc
					nm := mask | eventMask[combo]
					next[nst][nm] = (next[nst][nm] + cur) % mod
				}
			}
		}
		dp = next
	}

	ans := 0
	full := 64 - 1
	for st := 0; st < 8; st++ {
		ans = (ans + dp[st][full]) % mod
	}
	fmt.Println(ans)
}
