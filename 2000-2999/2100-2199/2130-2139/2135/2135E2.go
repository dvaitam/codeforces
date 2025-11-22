package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

// Precomputed linear recurrence (length 8) for the number of almost-palindrome
// binary strings. The sequence for n = 1.. is:
// 2, 2, 4, 8, 12, 26, 44, 86, ...
// Recurrence was derived via Berlekamp–Massey on the initial terms:
// a_n = sum_{i=1..8} coef[i-1] * a_{n-i} (mod mod)
var coef = []int64{
	426730412, 478501144, 629276108, 711585746,
	21639016, 531726114, 489204290, 957002268,
}

var initVals = []int64{
	0, // dummy for 1-based indexing
	2, 2, 4, 8, 12, 26, 44, 86,
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	ns := make([]int, T)
	maxN := 0
	for i := 0; i < T; i++ {
		fmt.Fscan(in, &ns[i])
		if ns[i] > maxN {
			maxN = ns[i]
		}
	}

	// Precompute answers up to maxN using the recurrence.
	ans := make([]int64, maxN+1) // 1-based
	limit := len(initVals) - 1
	for i := 1; i <= maxN && i <= limit; i++ {
		ans[i] = initVals[i]
	}
	for n := limit + 1; n <= maxN; n++ {
		var v int64
		for k := 1; k <= 8; k++ {
			v = (v + coef[k-1]*ans[n-k]) % mod
		}
		ans[n] = v
	}

	for i := 0; i < T; i++ {
		if i > 0 {
			fmt.Fprint(out, "\n")
		}
		fmt.Fprint(out, ans[ns[i]])
	}
}
