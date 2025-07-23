package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func mulDiv(a, b, div uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	q, _ := bits.Div64(hi, lo, div)
	return q
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k uint64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(in, &x)
		arr[n-1-i] = x // reverse order for convenience
	}

	limit := k
	has := func(t uint64) bool {
		var sum uint64
		comb := uint64(1) // C(t-1,0)
		for j := 0; j < n; j++ {
			if j > 0 {
				if t == 0 {
					break // further combinations are zero
				}
				comb = mulDiv(comb, t+uint64(j-1), uint64(j))
				if comb > limit {
					comb = limit + 1
				}
			}
			if arr[j] != 0 {
				if comb > limit/arr[j] {
					return true
				}
				sum += comb * arr[j]
				if sum >= limit {
					return true
				}
			}
			if comb == 0 {
				break
			}
		}
		return sum >= limit
	}

	lo, hi := uint64(0), k
	for lo < hi {
		mid := (lo + hi) / 2
		if has(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Fprintln(out, lo)
}
