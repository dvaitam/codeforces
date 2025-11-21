package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	dp := make([][2]int64, n)
	dp[0][0] = 1
	dp[0][1] = 1

	transform := func(bits [2]int64, prevDir, curDir int, target byte) int64 {
		total := int64(0)
		for prevPrevDir := 0; prevPrevDir < 2; prevPrevDir++ {
			if dp[prevDir][prevPrevDir] == 0 {
				continue
			}
			expected := nextState(prevPrevDir, prevDir)
			if expected == target {
				total += dp[prevDir][prevPrevDir]
			}
		}
		return total
	}

	for i := 1; i < n; i++ {
		var nextDP [2]int64
		for cur := 0; cur < 2; cur++ {
			for prev := 0; prev < 2; prev++ {
				expected := nextState(prev, cur)
				if expected == s[i] {
					nextDP[cur] += dp[i-1][prev]
				}
			}
		}
		dp[i] = nextDP
	}
	var ans int64
	for last := 0; last < 2; last++ {
		for prev := 0; prev < 2; prev++ {
			if dp[n-1][prev] == 0 {
				continue
			}
			expected := nextState(prev, last)
			if expected == s[0] && nextState(last, mapIndex(0)) == s[(n)%n] {
				ans += dp[n-1][prev]
			}
		}
	}
	fmt.Fprintln(out, ans)
}

func nextState(left, right int) byte {
	if left == 0 && right == 1 {
		return 'B'
	}
	if left == 1 && right == 0 {
		return 'A'
	}
	if left == 0 && right == 0 {
		return 'A'
	}
	return 'B'
}

func mapIndex(i int) int {
	return i % 2
}
