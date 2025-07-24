package main

import (
	"bufio"
	"fmt"
	"os"
)

func cmp(a, b byte) int8 {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

func solveGame(s string) string {
	n := len(s)
	dp := make([][]int8, n)
	for i := range dp {
		dp[i] = make([]int8, n)
	}
	for i := 0; i+1 < n; i++ {
		if s[i] == s[i+1] {
			dp[i][i+1] = 0
		} else {
			dp[i][i+1] = 1
		}
	}
	for length := 4; length <= n; length += 2 {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			// Alice takes left
			a1 := dp[l+2][r]
			if a1 == 0 {
				a1 = cmp(s[l], s[l+1])
			}
			a2 := dp[l+1][r-1]
			if a2 == 0 {
				a2 = cmp(s[l], s[r])
			}
			if a1 > a2 {
				a1 = a2
			}
			// Alice takes right
			b1 := dp[l][r-2]
			if b1 == 0 {
				b1 = cmp(s[r], s[r-1])
			}
			b2 := dp[l+1][r-1]
			if b2 == 0 {
				b2 = cmp(s[r], s[l])
			}
			if b1 > b2 {
				b1 = b2
			}
			if a1 < b1 {
				dp[l][r] = b1
			} else {
				dp[l][r] = a1
			}
		}
	}
	res := dp[0][n-1]
	if res > 0 {
		return "Alice"
	} else if res < 0 {
		return "Bob"
	}
	return "Draw"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, solveGame(s))
	}
}
