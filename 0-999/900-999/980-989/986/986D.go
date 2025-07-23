package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// maxProductForCost returns the maximum product achievable with total cost S
// using factors 2 (cost 2) and 3 (cost 3). It assumes S>0.
func maxProductForCost(S int64) *big.Int {
	var a, b int64
	r := S % 3
	switch r {
	case 0:
		b = S / 3
	case 1:
		if S < 4 {
			a = S / 2
		} else {
			b = (S - 4) / 3
			a = 2
		}
	case 2:
		if S < 2 {
			return big.NewInt(1)
		}
		b = (S - 2) / 3
		a = 1
	}
	res := new(big.Int).Exp(big.NewInt(3), big.NewInt(b), nil)
	if a > 0 {
		var t big.Int
		t.Exp(big.NewInt(2), big.NewInt(a), nil)
		res.Mul(res, &t)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	n := new(big.Int)
	n.SetString(s, 10)

	if n.Cmp(big.NewInt(1)) == 0 {
		fmt.Println(1)
		return
	}

	hi := int64(n.BitLen()*2 + 10)
	lo := int64(1)
	for lo < hi {
		mid := (lo + hi) / 2
		prod := maxProductForCost(mid)
		if prod.Cmp(n) >= 0 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Println(lo)
}
