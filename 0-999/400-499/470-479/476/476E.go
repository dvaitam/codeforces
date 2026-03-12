package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s, p string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &p)
	n := len(s)
	m := len(p)

	// For each start position i in s, compute:
	//   endPos[i] = minimum end position (exclusive) in s to match pattern p
	//   starting from position i (greedy match of p as subsequence)
	//   cost[i] = endPos[i] - i - m (number of extra characters consumed = removals within the match window)
	//   If no match possible, endPos[i] = -1
	endPos := make([]int, n)
	for i := 0; i < n; i++ {
		endPos[i] = -1
		pi := 0
		for j := i; j < n && pi < m; j++ {
			if s[j] == p[pi] {
				pi++
				if pi == m {
					endPos[i] = j + 1
				}
			}
		}
	}

	// dp[i][j] = max occurrences considering first i chars of s with j removals
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, n+1)
		for j := 0; j <= n; j++ {
			dp[i][j] = -1
		}
	}
	dp[0][0] = 0

	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			if dp[i][j] < 0 {
				continue
			}
			cur := dp[i][j]
			// Option 1: remove s[i]
			if j+1 <= n && dp[i+1][j+1] < cur {
				dp[i+1][j+1] = cur
			}
			// Option 2: keep s[i] but don't start a match here
			if dp[i+1][j] < cur {
				dp[i+1][j] = cur
			}
			// Option 3: start matching pattern at position i
			if endPos[i] >= 0 {
				e := endPos[i]
				removals := (e - i) - m // chars consumed but not matched
				newJ := j + removals
				if newJ <= n && dp[e][newJ] < cur+1 {
					dp[e][newJ] = cur + 1
				}
			}
		}
	}

	ans := make([]int, n+1)
	for x := 0; x <= n; x++ {
		if dp[n][x] > 0 {
			ans[x] = dp[n][x]
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprintf(writer, "%d", v)
	}
	writer.WriteByte('\n')
}
