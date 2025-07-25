package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program solves the problem described in problemC.txt (Sweets Eating).
// It computes the minimal possible sugar penalty for eating exactly k sweets
// for all k from 1..n when at most m sweets can be eaten per day.
// The optimal strategy is to sort the sweets by sugar concentration and then
// use dynamic programming with prefix sums.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}

	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if i <= m {
			dp[i] = prefix[i]
		} else {
			dp[i] = prefix[i] + dp[i-m]
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, dp[i])
	}
	fmt.Fprintln(out)
}
