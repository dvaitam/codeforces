package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	bonus := make([][]int64, n)
	for i := range bonus {
		bonus[i] = make([]int64, n)
	}
	for i := 0; i < k; i++ {
		var x, y int
		var c int64
		fmt.Fscan(in, &x, &y, &c)
		bonus[x-1][y-1] = c
	}

	maxMask := 1 << n
	dp := make([][]int64, maxMask)
	for i := range dp {
		dp[i] = make([]int64, n)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = a[i]
	}

	var ans int64
	for mask := 0; mask < maxMask; mask++ {
		cnt := bits.OnesCount(uint(mask))
		if cnt > m {
			continue
		}
		for last := 0; last < n; last++ {
			val := dp[mask][last]
			if val < 0 {
				continue
			}
			if cnt == m {
				if val > ans {
					ans = val
				}
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) != 0 {
					continue
				}
				newMask := mask | 1<<next
				nv := val + a[next] + bonus[last][next]
				if nv > dp[newMask][next] {
					dp[newMask][next] = nv
				}
			}
		}
	}

	fmt.Println(ans)
}
