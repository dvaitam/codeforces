package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const mod = 1000000007

// compute recursively the maximum value for subarray [l,r] (1-based)
// using memoization with big.Int
func maxValue(a []int64, l, r int, memo map[[2]int]*big.Int) *big.Int {
	if l == r {
		return big.NewInt(a[l-1])
	}
	key := [2]int{l, r}
	if v, ok := memo[key]; ok {
		return new(big.Int).Set(v)
	}
	best := big.NewInt(-1 << 60) // very small
	for k := l; k < r; k++ {
		left := maxValue(a, l, k, memo)
		right := maxValue(a, k+1, r, memo)
		// cand = left + 2*right
		cand := new(big.Int).Lsh(right, 1) // multiply by 2
		cand.Add(cand, left)
		if cand.Cmp(best) > 0 {
			best = cand
		}
	}
	memo[key] = new(big.Int).Set(best)
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		memo := make(map[[2]int]*big.Int)
		val := maxValue(a, l, r, memo)
		val.Mod(val, big.NewInt(mod))
		if val.Sign() < 0 {
			val.Add(val, big.NewInt(mod))
		}
		fmt.Fprintln(out, val)
	}
}
