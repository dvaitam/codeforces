package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const MOD int64 = 1000000007

// convert decimal string to base-p digits (least significant first)
func toBaseP(s string, p int64) []int {
	a := new(big.Int)
	a.SetString(s, 10)
	if a.Sign() == 0 {
		return []int{0}
	}
	pb := big.NewInt(p)
	digits := []int{}
	zero := big.NewInt(0)
	rem := new(big.Int)
	for a.Cmp(zero) > 0 {
		a.QuoRem(a, pb, rem)
		digits = append(digits, int(rem.Int64()))
	}
	return digits
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var p, alpha int64
	if _, err := fmt.Fscan(in, &p, &alpha); err != nil {
		return
	}
	var A string
	fmt.Fscan(in, &A)

	digits := toBaseP(A, p)
	L := len(digits)
	if alpha > int64(L) {
		fmt.Fprintln(out, 0)
		return
	}
	a := int(alpha)

	// dp[borrow][carry][less]
	dpPrev := make([][][]int64, 2)
	for i := range dpPrev {
		dpPrev[i] = make([][]int64, a+1)
		for j := range dpPrev[i] {
			dpPrev[i][j] = make([]int64, 2)
		}
	}
	dpPrev[0][0][0] = 1

	for pos := 0; pos < L; pos++ {
		ai := int64(digits[pos])
		dpNext := make([][][]int64, 2)
		for i := range dpNext {
			dpNext[i] = make([][]int64, a+1)
			for j := range dpNext[i] {
				dpNext[i][j] = make([]int64, 2)
			}
		}
		for borrow := 0; borrow <= 1; borrow++ {
			for c := 0; c <= a; c++ {
				for less := 0; less <= 1; less++ {
					val := dpPrev[borrow][c][less]
					if val == 0 {
						continue
					}
					if less == 0 {
						// choose n_i == ai
						if borrow == 0 {
							ways := (ai + 1) % MOD
							dpNext[0][c][0] = (dpNext[0][c][0] + val*ways) % MOD
						} else {
							// borrow == 1
							// k=n -> newBorrow=1
							c1 := c + 1
							if c1 > a {
								c1 = a
							}
							dpNext[1][c1][0] = (dpNext[1][c1][0] + val) % MOD
							// k<n -> newBorrow=0
							ways0 := ai % MOD
							dpNext[0][c][0] = (dpNext[0][c][0] + val*ways0) % MOD
						}
						// choose n_i < ai
						if ai > 0 {
							m := ai - 1
							if borrow == 0 {
								ways := ((m + 1) * (m + 2) / 2) % MOD
								dpNext[0][c][1] = (dpNext[0][c][1] + val*ways) % MOD
							} else {
								ways0 := (m * (m + 1) / 2) % MOD
								dpNext[0][c][1] = (dpNext[0][c][1] + val*ways0) % MOD
								ways1 := (m + 1) % MOD
								c1 := c + 1
								if c1 > a {
									c1 = a
								}
								dpNext[1][c1][1] = (dpNext[1][c1][1] + val*ways1) % MOD
							}
						}
					} else { // less == 1
						limit := int64(p) - 1
						if borrow == 0 {
							ways := ((limit + 1) * (limit + 2) / 2) % MOD
							dpNext[0][c][1] = (dpNext[0][c][1] + val*ways) % MOD
						} else {
							ways0 := (limit * (limit + 1) / 2) % MOD
							dpNext[0][c][1] = (dpNext[0][c][1] + val*ways0) % MOD
							ways1 := (limit + 1) % MOD
							c1 := c + 1
							if c1 > a {
								c1 = a
							}
							dpNext[1][c1][1] = (dpNext[1][c1][1] + val*ways1) % MOD
						}
					}
				}
			}
		}
		dpPrev = dpNext
	}

	ans := (dpPrev[0][a][0] + dpPrev[0][a][1]) % MOD
	fmt.Fprintln(out, ans)
}
