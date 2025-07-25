package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the wood log splitting problem described in problemC.txt.
// Starting with a log of weight 2^n, we may split any log into two halves.
// We want the minimum number of splits so that a subset of logs weighs exactly k.
// The optimal strategy recursively chooses a half that is closer to k and
// continues on that half. This leads to the recurrence:
//   f(n, k) = 0 if k == 0 or k == 2^n
//   f(n, k) = 1 + f(n-1, |k - 2^{n-1}|) otherwise.
// We iterate this relation until k becomes 0 or 2^n.

func minSplits(n int64, k int64) int64 {
	var steps int64
	for k > 0 && k < (int64(1)<<n) {
		half := int64(1) << (n - 1)
		if k > half {
			k -= half
		} else {
			k = half - k
		}
		steps++
		n--
	}
	return steps
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(in, &n, &k)
		fmt.Fprintln(out, minSplits(n, k))
	}
}
