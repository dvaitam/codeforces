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

	var t int
	fmt.Fscan(in, &t)
	pattern := []byte{'n', 'a', 'r', 'e', 'k'}

	const negInf = int64(-1 << 60)

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		strings := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &strings[i])
		}

		totals := make([]int, n)
		increments := make([][5]int, n)
		nextState := make([][5]int, n)

		for idx, s := range strings {
			bytes := []byte(s)
			cnt := 0
			for _, ch := range bytes {
				if ch == 'n' || ch == 'a' || ch == 'r' || ch == 'e' || ch == 'k' {
					cnt++
				}
			}
			totals[idx] = cnt

			for start := 0; start < 5; start++ {
				state := start
				complete := 0
				for _, ch := range bytes {
					if ch == pattern[state] {
						state++
						if state == 5 {
							complete++
							state = 0
						}
					}
				}
				increments[idx][start] = complete
				nextState[idx][start] = state
			}
		}

		dp := [5]int64{}
		for i := 0; i < 5; i++ {
			dp[i] = negInf
		}
		dp[0] = 0

		for i := 0; i < n; i++ {
			var newDP [5]int64
			for s := 0; s < 5; s++ {
				newDP[s] = dp[s]
			}
			for s := 0; s < 5; s++ {
				if dp[s] == negInf {
					continue
				}
				inc := int64(increments[i][s])
				ns := nextState[i][s]
				delta := 10*inc - int64(totals[i])
				candidate := dp[s] + delta
				if candidate > newDP[ns] {
					newDP[ns] = candidate
				}
			}
			dp = newDP
		}

		ans := dp[0]
		for s := 1; s < 5; s++ {
			if dp[s] > ans {
				ans = dp[s]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
