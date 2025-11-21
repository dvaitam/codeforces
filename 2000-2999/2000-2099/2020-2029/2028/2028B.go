package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

var (
	one = big.NewInt(1)
	two = big.NewInt(2)
)

func parseBig(s string) *big.Int {
	val := new(big.Int)
	if _, ok := val.SetString(s, 10); !ok {
		panic("invalid number")
	}
	return val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var ns, bs, cs string
		fmt.Fscan(in, &ns, &bs, &cs)
		n := parseBig(ns)
		b := parseBig(bs)
		c := parseBig(cs)

		if b.Sign() == 0 {
			nGreaterOne := n.Cmp(one) == 1
			threshold := new(big.Int).Add(c, two)
			if nGreaterOne && n.Cmp(threshold) == 1 {
				fmt.Fprintln(out, -1)
				continue
			}
			ans := new(big.Int).Set(n)
			if c.Cmp(n) == -1 {
				ans.Sub(ans, one)
			}
			fmt.Fprintln(out, ans.String())
			continue
		}

		maxVal := new(big.Int).Sub(new(big.Int).Set(n), one)
		count := big.NewInt(0)
		if c.Cmp(maxVal) <= 0 {
			tRange := new(big.Int).Sub(maxVal, c)
			quotient := new(big.Int).Div(tRange, b)
			count = quotient.Add(quotient, one)
			if count.Cmp(n) == 1 {
				count.Set(n)
			}
		}
		ans := new(big.Int).Sub(new(big.Int).Set(n), count)
		fmt.Fprintln(out, ans.String())
	}
}
