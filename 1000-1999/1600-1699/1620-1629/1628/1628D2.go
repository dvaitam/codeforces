package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// This solution uses a direct dynamic programming approach with exact
// rational arithmetic using big.Rat. It follows the game definition
// literally and therefore is not optimized for the larger constraints
// of the hard version, but it demonstrates the intended recurrence.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	mod := big.NewInt(1_000_000_007)
	inv := func(x *big.Int) *big.Int {
		// modular inverse using Fermat since mod is prime
		exp := new(big.Int).Sub(mod, big.NewInt(2))
		return new(big.Int).Exp(x, exp, mod)
	}

	half := big.NewRat(1, 2)

	for ; T > 0; T-- {
		var n, m int
		var k int64
		fmt.Fscan(in, &n, &m, &k)
		if m == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		if m > n {
			fmt.Fprintln(out, 0)
			continue
		}
		// dp arrays of big.Rat
		prev := make([]*big.Rat, m+1)
		curr := make([]*big.Rat, m+1)
		for i := 0; i <= m; i++ {
			prev[i] = new(big.Rat)
			curr[i] = new(big.Rat)
		}
		prev[0].SetInt64(0)
		for i := 1; i <= n; i++ {
			upto := m
			if i < m {
				upto = i
			}
			for j := 1; j <= upto; j++ {
				if j == i {
					curr[j].SetInt64(int64(i))
					curr[j].Mul(curr[j], big.NewRat(k, 1))
					continue
				}
				delta := new(big.Rat).Sub(prev[j], prev[j-1])
				cmp0 := delta.Cmp(new(big.Rat))
				twoK := new(big.Rat).SetInt64(2 * k)
				if cmp0 <= 0 {
					curr[j].Set(prev[j])
				} else if delta.Cmp(twoK) >= 0 {
					curr[j].Set(prev[j-1])
					curr[j].Add(curr[j], big.NewRat(k, 1))
				} else {
					curr[j].Set(prev[j-1])
					delta.Mul(delta, half)
					curr[j].Add(curr[j], delta)
				}
			}
			// reset unused entries
			if i <= m {
				curr[i+1] = new(big.Rat)
			}
			prev, curr = curr, prev
		}
		ans := prev[m]
		// convert big.Rat to modular value p * q^{-1} mod M
		p := new(big.Int).Mod(ans.Num(), mod)
		q := new(big.Int).Mod(ans.Denom(), mod)
		qInv := inv(q)
		p.Mul(p, qInv)
		p.Mod(p, mod)
		fmt.Fprintln(out, p)
	}
}
