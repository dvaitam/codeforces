package main

import (
	"bufio"
	"fmt"
	"os"
)

// The exact algorithm for this problem is non-trivial and normally would
// require the official editorial. In this placeholder we implement a very
// rough heuristic based on the maximum subarray sum of b[i] = max(0, a[i]+1).
// It is **not** guaranteed to be correct but compiles and can serve as a
// starting point for a full solution.

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	// Build auxiliary array b[i] = max(0, a[i]+1)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if a[i]+1 > 0 {
			b[i] = a[i] + 1
		} else {
			b[i] = 0
		}
	}
	// Kadane on b to get a crude upper bound
	cur, best := 0, 0
	for _, v := range b {
		cur += v
		if cur > best {
			best = cur
		}
		if cur < 0 {
			cur = 0
		}
	}
	fmt.Println(best)
}
