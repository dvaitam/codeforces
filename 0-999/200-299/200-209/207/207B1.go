package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n+1) // 1-indexed by tank number
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i + 1
	}

	var total int64
	dp := make([]int, n)

	for msg := 0; msg < n; msg++ {
		for i := range dp {
			dp[i] = INF
		}
		dp[0] = 0 // first tank already has the message

		for i := 0; i < n; i++ {
			if dp[i] == INF {
				continue
			}
			for j := i + 1; j < n; j++ {
				receiver := order[j]
				if (i + 1) >= (j+1)-a[receiver] {
					if dp[j] > dp[i]+1 {
						dp[j] = dp[i] + 1
					}
				}
			}
		}

		total += int64(dp[n-1])

		// rotate: move last tank to the front
		last := order[n-1]
		copy(order[1:], order[:n-1])
		order[0] = last
	}

	fmt.Fprintln(out, total)
}
