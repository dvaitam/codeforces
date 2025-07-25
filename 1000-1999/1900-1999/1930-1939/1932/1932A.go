package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// We walk along a path of n cells, each either empty '.',
// containing a coin '@', or a thorn '*'. From the first cell we
// may move one or two cells to the right as long as the
// destination is not a thorn. The goal is to collect the maximum
// number of coins we can reach.
//
// A simple dynamic programming approach suffices since n <= 50.
// Let dp[i] be the maximum number of coins that can be collected
// when arriving at cell i (0-indexed) or -1 if cell i is
// unreachable. For every non-thorn cell i we look at dp[i-1] and
// dp[i-2] (when i >= 2) and choose the best reachable previous
// cell. If cell i contains a coin we add one to the value. The
// answer for each test case is the maximum dp[i] over all cells
// because we may stop anywhere we can reach.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		dp := make([]int, n)
		for i := range dp {
			dp[i] = -1
		}
		dp[0] = 0 // first cell is guaranteed empty
		for i := 1; i < n; i++ {
			if s[i] == '*' {
				continue
			}
			best := -1
			if dp[i-1] != -1 && dp[i-1] > best {
				best = dp[i-1]
			}
			if i >= 2 && dp[i-2] != -1 && dp[i-2] > best {
				best = dp[i-2]
			}
			if best == -1 {
				continue
			}
			if s[i] == '@' {
				best++
			}
			dp[i] = best
		}
		ans := 0
		for _, v := range dp {
			if v > ans {
				ans = v
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
