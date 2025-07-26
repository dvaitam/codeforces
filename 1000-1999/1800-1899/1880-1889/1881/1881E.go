package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution implements a dynamic programming approach to find the
// longest subsequence which already has the "beautiful" form described in
// the statement. Each block of the subsequence consists of the block length
// followed by that many elements. When starting a block at position i with
// length a[i], taking the next a[i] elements greedily minimises the ending
// index of the block and thus is always optimal. Therefore the problem
// reduces to selecting non-overlapping segments [i, i+a[i]] to maximise
// their total length.
//
// dp[i] holds the maximum length achievable using the suffix starting at i
// (1-indexed). For each i we either skip a[i] or take the block starting at
// i if i+a[i] <= n. The answer is n - dp[1].
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		dp := make([]int, n+2)
		for i := n; i >= 1; i-- {
			dp[i] = dp[i+1]
			l := a[i]
			if i+l <= n {
				val := l + 1 + dp[i+l+1]
				if val > dp[i] {
					dp[i] = val
				}
			}
		}

		fmt.Fprintln(out, n-dp[1])
	}
}
