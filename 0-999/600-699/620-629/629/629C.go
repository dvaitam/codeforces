package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	delta := 0
	minPref := 0
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < minPref {
			minPref = bal
		}
	}
	delta = bal

	maxLen := n - m
	dp := make([][]int64, maxLen+1)
	for i := range dp {
		dp[i] = make([]int64, maxLen+1)
	}
	dp[0][0] = 1
	for i := 0; i < maxLen; i++ {
		for j := 0; j <= maxLen; j++ {
			v := dp[i][j]
			if v == 0 {
				continue
			}
			if j+1 <= maxLen {
				dp[i+1][j+1] = (dp[i+1][j+1] + v) % mod
			}
			if j > 0 {
				dp[i+1][j-1] = (dp[i+1][j-1] + v) % mod
			}
		}
	}

	ans := int64(0)
	for l := 0; l <= maxLen; l++ {
		qlen := maxLen - l
		for j := 0; j <= maxLen; j++ {
			v := dp[l][j]
			if v == 0 {
				continue
			}
			if j+minPref < 0 {
				continue
			}
			y := j + delta
			if y < 0 || y > maxLen {
				continue
			}
			add := dp[qlen][y]
			if add == 0 {
				continue
			}
			ans = (ans + v*add) % mod
		}
	}

	fmt.Fprintln(out, ans)
}
