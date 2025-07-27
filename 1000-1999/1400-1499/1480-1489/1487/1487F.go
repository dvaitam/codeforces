package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	n := new(big.Int)
	n.SetString(s, 10)
	if n.Sign() == 0 {
		fmt.Println(0)
		return
	}

	// precompute repunits of lengths up to 60
	reps := make([]*big.Int, 0, 60)
	one := big.NewInt(1)
	cur := big.NewInt(0)
	for i := 0; i < 60; i++ {
		cur = new(big.Int).Add(new(big.Int).Mul(cur, big.NewInt(10)), one)
		reps = append(reps, new(big.Int).Set(cur))
	}

	cost := 0
	abs := new(big.Int).Abs(n)
	zero := big.NewInt(0)

	for abs.Cmp(zero) != 0 {
		bestDiff := new(big.Int).Set(abs)
		bestK := 1
		bestNew := new(big.Int)
		for k, r := range reps {
			// subtract r
			tmp := new(big.Int).Sub(abs, r)
			if tmp.Sign() < 0 {
				tmp.Abs(tmp)
			}
			if tmp.Cmp(bestDiff) < 0 || (tmp.Cmp(bestDiff) == 0 && k+1 < bestK) {
				bestDiff.Set(tmp)
				bestK = k + 1
				bestNew.Sub(abs, r)
			}
			// add r
			tmp = new(big.Int).Add(abs, r)
			if tmp.Sign() < 0 {
				tmp.Abs(tmp)
			}
			if tmp.Cmp(bestDiff) < 0 || (tmp.Cmp(bestDiff) == 0 && k+1 < bestK) {
				bestDiff.Set(tmp)
				bestK = k + 1
				bestNew.Add(abs, r)
			}
		}
		abs = bestNew.Abs(bestNew)
		cost += bestK
	}

	fmt.Println(cost)
}
