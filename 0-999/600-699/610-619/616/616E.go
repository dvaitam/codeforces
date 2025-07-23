package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const MOD int64 = 1e9 + 7

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	limit := n
	if m < limit {
		limit = m
	}
	modBig := big.NewInt(MOD)
	res := big.NewInt(0)
	l := int64(1)
	for l <= limit {
		q := n / l
		r := n / q
		if r > limit {
			r = limit
		}
		count := (r - l + 1) % MOD
		term1 := new(big.Int).Mul(big.NewInt(n%MOD), big.NewInt(count))
		term1.Mod(term1, modBig)

		sum := new(big.Int).Mul(big.NewInt(l+r), big.NewInt(r-l+1))
		sum.Div(sum, big.NewInt(2))
		sum.Mod(sum, modBig)

		term2 := new(big.Int).Mul(big.NewInt(q%MOD), sum)
		term2.Mod(term2, modBig)

		term1.Sub(term1, term2)
		term1.Mod(term1, modBig)
		res.Add(res, term1)
		res.Mod(res, modBig)

		l = r + 1
	}
	if m > n {
		extra := new(big.Int).Mul(big.NewInt((m-n)%MOD), big.NewInt(n%MOD))
		extra.Mod(extra, modBig)
		res.Add(res, extra)
		res.Mod(res, modBig)
	}
	fmt.Fprintln(out, res.Int64())
}
