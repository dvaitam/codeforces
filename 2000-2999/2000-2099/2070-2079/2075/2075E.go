package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// fast power of 2 modulo mod.
func pow2(e int64) int64 {
	res := int64(1)
	base := int64(2)
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

// counts quadruples (a,b,c,d) with a,b in [0,A], c,d in [0,B], and a^b^c^d == 0.
func countXorZero(A, B int64) int64 {
	maxBit := 0
	for (1 << maxBit) <= int(A|B) {
		maxBit++
	}
	if maxBit == 0 {
		maxBit = 1
	}

	dp := [16]int64{}
	dp[15] = 1 // all tight at start

	for bit := maxBit - 1; bit >= 0; bit-- {
		limA := int((A >> bit) & 1)
		limB := int((B >> bit) & 1)
		var ndp [16]int64
		for state := 0; state < 16; state++ {
			cur := dp[state]
			if cur == 0 {
				continue
			}
			tightA1 := (state>>0)&1 == 1
			tightB1 := (state>>1)&1 == 1
			tightA2 := (state>>2)&1 == 1
			tightB2 := (state>>3)&1 == 1

			for ba := 0; ba <= 1; ba++ {
				if tightA1 && ba > limA {
					continue
				}
				naTight := tightA1 && ba == limA
				for bb := 0; bb <= 1; bb++ {
					if tightB1 && bb > limA { // b bound is A
						continue
					}
					nbTight := tightB1 && bb == limA
					for bc := 0; bc <= 1; bc++ {
						if tightA2 && bc > limB { // c bound is B
							continue
						}
						ncTight := tightA2 && bc == limB
						bd := ba ^ bb ^ bc
						if bd > 1 {
							continue
						}
						if tightB2 && bd > limB { // d bound is B
							continue
						}
						ndTight := tightB2 && bd == limB

						nstate := 0
						if naTight {
							nstate |= 1 << 0
						}
						if nbTight {
							nstate |= 1 << 1
						}
						if ncTight {
							nstate |= 1 << 2
						}
						if ndTight {
							nstate |= 1 << 3
						}
						ndp[nstate] = (ndp[nstate] + cur) % mod
					}
				}
			}
		}
		dp = ndp
	}

	sum := int64(0)
	for _, v := range dp {
		sum = (sum + v) % mod
	}
	return sum
}

func solveCase(n, m, A, B int64) int64 {
	oneVal := ((A + 1) % mod) * ((B + 1) % mod) % mod

	// total quadruples where xor zero
	total := countXorZero(A, B)
	// gA[0] = A+1, gB[0] = B+1
	nonZero := (total - (A+1)%mod*((B+1)%mod)%mod) % mod
	if nonZero < 0 {
		nonZero += mod
	}

	p2n := pow2(n)
	p2m := pow2(m)
	factor := (p2n * p2m) % mod
	factor = (factor - 4) % mod
	if factor < 0 {
		factor += mod
	}

	ans := (oneVal + factor*nonZero%mod) % mod
	return ans
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
		var n, m, A, B int64
		fmt.Fscan(in, &n, &m, &A, &B)
		fmt.Fprintln(out, solveCase(n, m, A, B))
	}
}
