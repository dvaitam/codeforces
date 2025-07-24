package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func countWaysClose(s string) int {
	open := 0
	close := 0
	// positions of ')' indices (1-based)
	var p []int
	idx := 0
	for _, c := range s {
		if c == '(' {
			open++
		} else {
			close++
			idx++
			if close > open {
				p = append(p, idx)
				open = close
			}
		}
	}
	m := idx
	r := len(p)
	if r == 0 {
		return 1
	}
	dp := make([]int, r+1)
	dp[0] = 1
	for j := 1; j <= m; j++ {
		limit := r
		if j < limit {
			limit = j
		}
		for t := limit; t >= 1; t-- {
			if j <= p[t-1] {
				dp[t] += dp[t-1]
				if dp[t] >= MOD {
					dp[t] -= MOD
				}
			}
		}
	}
	return dp[r]
}

func countWaysOpen(s string) int {
	open := 0
	close := 0
	var q []int
	idx := 0
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c == ')' {
			close++
		} else {
			open++
			idx++
			if open > close {
				q = append(q, idx)
				close = open
			}
		}
	}
	m := idx
	r := len(q)
	if r == 0 {
		return 1
	}
	dp := make([]int, r+1)
	dp[0] = 1
	for j := 1; j <= m; j++ {
		limit := r
		if j < limit {
			limit = j
		}
		for t := limit; t >= 1; t-- {
			if j <= q[t-1] {
				dp[t] += dp[t-1]
				if dp[t] >= MOD {
					dp[t] -= MOD
				}
			}
		}
	}
	return dp[r]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	waysClose := countWaysClose(s)
	waysOpen := countWaysOpen(s)
	ans := int(int64(waysClose) * int64(waysOpen) % int64(MOD))
	fmt.Println(ans)
}
