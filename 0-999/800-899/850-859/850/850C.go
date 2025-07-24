package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// compute Grundy number for bitmask representing exponents for a single prime.
// bit i (i>=1) is set if there exists at least one number with exponent i.
var memo = map[int]int{}

func grundy(mask int) int {
	if mask == 0 {
		return 0
	}
	if v, ok := memo[mask]; ok {
		return v
	}
	seen := map[int]bool{}
	maxBit := bits.Len(uint(mask))
	for k := 1; k <= maxBit; k++ {
		if mask>>k == 0 {
			break
		}
		nm := (mask & ((1 << k) - 1)) | ((mask >> k) &^ 1)
		g := grundy(nm)
		seen[g] = true
	}
	g := 0
	for {
		if !seen[g] {
			memo[mask] = g
			return g
		}
		g++
	}
}

func factorize(x int) map[int]int {
	res := map[int]int{}
	tmp := x
	for p := 2; p*p <= tmp; p++ {
		if tmp%p == 0 {
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			res[p] = cnt
		}
	}
	if tmp > 1 {
		res[tmp]++
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	primeMask := map[int]int{}
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		fac := factorize(x)
		for p, e := range fac {
			primeMask[p] |= 1 << e
		}
	}
	nim := 0
	for _, mask := range primeMask {
		nim ^= grundy(mask)
	}
	if nim != 0 {
		fmt.Println("Mojtaba")
	} else {
		fmt.Println("Arpa")
	}
}
