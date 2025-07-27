package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var ls, rs string
	fmt.Fscan(in, &ls)
	fmt.Fscan(in, &rs)
	l := new(big.Int)
	r := new(big.Int)
	l.SetString(ls, 2)
	r.SetString(rs, 2)

	k := r.BitLen() - 1
	candidate := new(big.Int)
	found := false
	for ; k >= 0; k-- {
		two := new(big.Int).Lsh(big.NewInt(1), uint(k))
		if r.Cmp(two) < 0 {
			continue
		}
		minusOne := new(big.Int).Sub(two, big.NewInt(1))
		if l.Cmp(minusOne) <= 0 {
			candidate.Lsh(big.NewInt(1), uint(k+1))
			candidate.Sub(candidate, big.NewInt(1))
			found = true
			break
		}
	}
	if !found {
		temp := new(big.Int).Sub(r, big.NewInt(1))
		candidate.Or(r, temp)
	}

	if r.Cmp(candidate) > 0 {
		fmt.Println(r.Text(2))
	} else {
		fmt.Println(candidate.Text(2))
	}
}
