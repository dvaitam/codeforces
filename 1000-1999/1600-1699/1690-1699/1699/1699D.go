package main

// Solution for problem D as described in problemD.txt.
// We are given an array and may repeatedly delete adjacent pairs of
// different elements. The task is to obtain the longest possible array
// consisting of equal elements only.
//
// Key observation: a segment can be completely removed if and only if its
// length is even and the frequency of its most common value does not
// exceed half of the segment length. Using this, we precompute for every
// interval whether it can vanish. Then dynamic programming is applied:
// dp[i] is the maximum length of a sequence of the same value that ends at
// position i such that everything before i in the array can be deleted.
// For each position we transition from previous positions with the same
// value if the segment between them is removable. The answer is the best
// dp[i] for which the suffix after i is also removable.

import (
	"bufio"
	"fmt"
	"os"
)

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

		can := make([][]bool, n+2)
		for i := range can {
			can[i] = make([]bool, n+2)
		}
		for i := 1; i <= n+1; i++ {
			can[i][i-1] = true
		}

		for l := 1; l <= n; l++ {
			freq := make([]int, n+1)
			maxf := 0
			for r := l; r <= n; r++ {
				x := a[r]
				freq[x]++
				if freq[x] > maxf {
					maxf = freq[x]
				}
				length := r - l + 1
				if length%2 == 0 && maxf <= length/2 {
					can[l][r] = true
				}
			}
		}

		dp := make([]int, n+1)
		ans := 0
		for i := 1; i <= n; i++ {
			if can[1][i-1] {
				dp[i] = 1
			}
			for j := 1; j < i; j++ {
				if a[i] == a[j] && dp[j] > 0 && can[j+1][i-1] {
					if dp[j]+1 > dp[i] {
						dp[i] = dp[j] + 1
					}
				}
			}
			if can[i+1][n] && dp[i] > ans {
				ans = dp[i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
