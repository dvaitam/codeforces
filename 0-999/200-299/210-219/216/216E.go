package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var k, b int64
	var n int
	if _, err := fmt.Fscan(reader, &k, &b, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// prefix sums
	P := make([]int64, n+1)
	for i := 0; i < n; i++ {
		P[i+1] = P[i] + a[i]
	}
	// digital root logic uses M = k-1
	M := k - 1
	var result int64
	// b == 0: only zero-sum substrings
	if b == 0 {
		cnt := make(map[int64]int64)
		for _, v := range P {
			cnt[v]++
		}
		for _, c := range cnt {
			if c > 1 {
				result += c * (c - 1) / 2
			}
		}
		fmt.Println(result)
		return
	}
	// b == M (i.e., k-1): substrings with sum mod M == 0 and sum > 0
	if b == M {
		cntMod := make(map[int64]int64)
		cntExact := make(map[int64]int64)
		var totalMod, totalExact int64
		// initialize with P[0]
		cntMod[P[0]%M] = 1
		cntExact[P[0]] = 1
		for j := 1; j <= n; j++ {
			r := P[j] % M
			// total substrings with mod == 0: P[i] mod M == r
			totalMod += cntMod[r]
			// substrings with exact zero sum: P[i] == P[j]
			totalExact += cntExact[P[j]]
			cntMod[r]++
			cntExact[P[j]]++
		}
		result = totalMod - totalExact
		fmt.Println(result)
		return
	}
	// general b in 1..M-1: substrings with sum mod M == b
	cntMod := make(map[int64]int64)
	// initialize with P[0]
	cntMod[P[0]%M] = 1
	for j := 1; j <= n; j++ {
		r := P[j] % M
		// need P[i] mod M == (r - b mod M)
		need := r - b
		need %= M
		if need < 0 {
			need += M
		}
		if c, ok := cntMod[need]; ok {
			result += c
		}
		cntMod[r]++
	}
	fmt.Println(result)
}
