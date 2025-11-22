package main

import (
	"bufio"
	"fmt"
	"os"
)

type parent struct {
	prevLow  int
	prevHigh int
	prevEq   int
	bits     int // bits chosen for a,b,c at this position (bit0 -> a)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	for ; T > 0; T-- {
		var l, r int
		fmt.Fscan(in, &l, &r)

		const bits = 30
		var dp [bits + 1][8][8][8]int64
		for i := 0; i <= bits; i++ {
			for a := 0; a < 8; a++ {
				for b := 0; b < 8; b++ {
					for c := 0; c < 8; c++ {
						dp[i][a][b][c] = -1
					}
				}
			}
		}
		var par [bits][8][8][8]parent

		dp[bits][7][7][7] = 0 // all tight, all equal

		for idx := bits - 1; idx >= 0; idx-- {
			lb := (l >> idx) & 1
			rb := (r >> idx) & 1
			for low := 0; low < 8; low++ {
				for high := 0; high < 8; high++ {
					for eq := 0; eq < 8; eq++ {
						cur := dp[idx+1][low][high][eq]
						if cur < 0 {
							continue
						}
						for mask := 0; mask < 8; mask++ { // bits for a,b,c
							b0 := mask & 1
							b1 := (mask >> 1) & 1
							b2 := (mask >> 2) & 1
							// validate against bounds
							ok := true
							nl, nh := 0, 0
							for j, bj := range []int{b0, b1, b2} {
								tl := (low >> j) & 1
								th := (high >> j) & 1
								if tl == 1 && bj < lb {
									ok = false
									break
								}
								if th == 1 && bj > rb {
									ok = false
									break
								}
								if tl == 1 && bj == lb {
									nl |= 1 << j
								}
								if th == 1 && bj == rb {
									nh |= 1 << j
								}
							}
							if !ok {
								continue
							}

							neq := 0
							// eq bits: 0->ab,1->ac,2->bc
							if (eq&1) == 1 && b0 == b1 {
								neq |= 1
							}
							if (eq&2) == 2 && b0 == b2 {
								neq |= 2
							}
							if (eq&4) == 4 && b1 == b2 {
								neq |= 4
							}

							diffPairs := (b0 ^ b1) + (b0 ^ b2) + (b1 ^ b2)
							val := cur + int64(diffPairs<<idx)
							if val > dp[idx][nl][nh][neq] {
								dp[idx][nl][nh][neq] = val
								par[idx][nl][nh][neq] = parent{
									prevLow:  low,
									prevHigh: high,
									prevEq:   eq,
									bits:     mask,
								}
							}
						}
					}
				}
			}
		}

		bestVal := int64(-1)
		bestLow, bestHigh := 0, 0
		// eq mask must be 0 (all pairs different)
		for low := 0; low < 8; low++ {
			for high := 0; high < 8; high++ {
				if dp[0][low][high][0] > bestVal {
					bestVal = dp[0][low][high][0]
					bestLow, bestHigh = low, high
				}
			}
		}

		// reconstruct numbers
		aBits := 0
		bBits := 0
		cBits := 0
		low, high, eq := bestLow, bestHigh, 0
		for idx := 0; idx < bits; idx++ {
			p := par[idx][low][high][eq]
			mask := p.bits
			if (mask & 1) == 1 {
				aBits |= 1 << idx
			}
			if (mask & 2) == 2 {
				bBits |= 1 << idx
			}
			if (mask & 4) == 4 {
				cBits |= 1 << idx
			}
			low, high, eq = p.prevLow, p.prevHigh, p.prevEq
		}

		fmt.Fprintf(out, "%d %d %d\n", aBits, bBits, cBits)
	}
	out.Flush()
}
