package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const MOD = 1000000007
const LIMIT = 200000 // maximum iterations for enumeration

// checkCandidate checks whether x belongs to all given progressions
func checkCandidate(x *big.Int, a, b []int64) bool {
	tmp := new(big.Int)
	for i := range a {
		// x must be divisible by a[i]
		tmp.Mod(x, big.NewInt(a[i]))
		if tmp.Sign() != 0 {
			return false
		}
		tmp.Div(x, big.NewInt(a[i]))
		if b[i] == 1 {
			if tmp.Cmp(big.NewInt(1)) != 0 {
				return false
			}
			continue
		}
		bi := big.NewInt(b[i])
		for tmp.Mod(tmp, bi).Sign() == 0 {
			tmp.Div(tmp, bi)
		}
		if tmp.Cmp(big.NewInt(1)) != 0 {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
	}

	// handle constant progressions first
	constVal := int64(-1)
	for i := 0; i < n; i++ {
		if b[i] == 1 {
			if constVal == -1 {
				constVal = a[i]
			} else if constVal != a[i] {
				fmt.Println(-1)
				return
			}
		}
	}
	if constVal != -1 {
		// verify this value exists in all progressions
		cand := big.NewInt(constVal)
		if checkCandidate(cand, a, b) {
			modRes := new(big.Int).Mod(cand, big.NewInt(MOD))
			fmt.Println(modRes.String())
		} else {
			fmt.Println(-1)
		}
		return
	}

	// choose progression with smallest ratio as base for enumeration
	base := 0
	for i := 1; i < n; i++ {
		if b[i] < b[base] || (b[i] == b[base] && a[i] < a[base]) {
			base = i
		}
	}

	cand := big.NewInt(a[base])
	ratio := big.NewInt(b[base])

	for k := 0; k < LIMIT; k++ {
		if checkCandidate(cand, a, b) {
			modRes := new(big.Int).Mod(cand, big.NewInt(MOD))
			fmt.Println(modRes.String())
			return
		}
		if b[base] == 1 {
			break
		}
		cand.Mul(cand, ratio)
	}

	fmt.Println(-1)
}
