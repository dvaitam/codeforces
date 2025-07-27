package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problem E from contest 1550.
// We are given a string s of length n containing the first k letters and '?'.
// We may replace every '?' with a letter from ['a'..'a'+k-1].
// For each letter i we define f_i as the length of the longest substring
// consisting solely of that letter. The value of the string is min_i f_i.
// The goal is to maximize this value after replacements.
//
// For a candidate length L we can check feasibility using dynamic programming.
// For each letter we precompute the earliest starting position of a segment of
// length L that can become that letter (only that letter or '?'). Using a DP
// over bitmasks of letters we greedily place non-overlapping segments of length
// L in increasing order. If we can place segments for all k letters the value L
// is achievable. Binary search on L gives the optimal answer.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	bytes := []byte(s)

	low, high := 0, n
	ans := 0
	for low <= high {
		mid := (low + high) / 2
		if possible(bytes, n, k, mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	fmt.Fprintln(out, ans)
}

func possible(s []byte, n, k, L int) bool {
	if L == 0 {
		return true
	}
	if k*L > n {
		return false
	}
	const INF = int(1e9)

	next := make([][]int, k)
	prefix := make([]int, n+1)

	for c := 0; c < k; c++ {
		// compute prefix of characters that cannot be letter c
		prefix[0] = 0
		target := byte('a' + c)
		for i := 0; i < n; i++ {
			if s[i] != '?' && s[i] != target {
				prefix[i+1] = prefix[i] + 1
			} else {
				prefix[i+1] = prefix[i]
			}
		}
		nxt := make([]int, n+1)
		nxt[n] = INF
		for i := n - 1; i >= 0; i-- {
			nxt[i] = nxt[i+1]
			if i+L <= n && prefix[i+L]-prefix[i] == 0 {
				nxt[i] = i
			}
		}
		next[c] = nxt
	}

	dp := make([]int, 1<<k)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	for mask := 0; mask < (1 << k); mask++ {
		pos := dp[mask]
		if pos > n {
			continue
		}
		for c := 0; c < k; c++ {
			if mask&(1<<c) != 0 {
				continue
			}
			start := next[c][pos]
			if start == INF {
				continue
			}
			end := start + L
			if end <= n && end < dp[mask|(1<<c)] {
				dp[mask|(1<<c)] = end
			}
		}
	}

	return dp[(1<<k)-1] <= n
}
