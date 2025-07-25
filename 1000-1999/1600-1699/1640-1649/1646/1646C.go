package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// This program solves the problem described in problemC.txt.
// We need to represent n as the sum of distinct powerful numbers,
// where powerful numbers are powers of two or factorials. We
// enumerate all subsets of factorials up to 14! (the largest <= 1e12)
// and for each subset compute how many additional powers of two are
// needed. The minimum over all subsets is the answer.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// precompute factorials up to 14!
	fac := []int64{1}
	for i := int64(2); ; i++ {
		next := fac[len(fac)-1] * i
		if next > 1_000_000_000_000 {
			break
		}
		fac = append(fac, next)
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		best := 60 // upper bound since n<2^40 and max factorials count<14
		m := len(fac)
		for mask := 0; mask < (1 << m); mask++ {
			sum := int64(0)
			cnt := 0
			for i := 0; i < m; i++ {
				if mask>>i&1 == 1 {
					sum += fac[i]
					cnt++
				}
			}
			if sum > n {
				continue
			}
			remainder := n - sum
			k := cnt + bits.OnesCount64(uint64(remainder))
			if k < best {
				best = k
			}
		}
		fmt.Fprintln(out, best)
	}
}
