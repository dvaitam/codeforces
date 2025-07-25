package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		if _, err := fmt.Fscan(in, &n, &k); err != nil {
			return
		}

		dp := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			dp[i] = make([]int64, k)
		}

		for a0 := 0; a0 < k; a0++ {
			for a1 := 0; a1 < k-a0; a1++ {
				w := (a0 + 1) * (a1 + 1)
				if w <= n {
					dp[w][a1]++
				}
			}
		}

		for s := 1; s <= n; s++ {
			for prev := 0; prev < k; prev++ {
				val := dp[s][prev]
				if val == 0 {
					continue
				}
				for nxt := 0; nxt < k-prev; nxt++ {
					w := (prev + 1) * (nxt + 1)
					ns := s + w
					if ns <= n {
						dp[ns][nxt] = (dp[ns][nxt] + val) % mod
					}
				}
			}
		}

		var res int64
		for last := 0; last < k; last++ {
			res = (res + dp[n][last]) % mod
		}
		fmt.Fprintln(out, res)
	}
}
