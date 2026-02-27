package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, T int
	fmt.Fscan(in, &n, &T)

	tasks := make([][]int, T+1)
	for i := 0; i < n; i++ {
		var t, q int
		fmt.Fscan(in, &t, &q)
		tasks[t] = append(tasks[t], q)
	}

	for t := 1; t <= T; t++ {
		sort.Sort(sort.Reverse(sort.IntSlice(tasks[t])))
	}

	prefix := make([][]int, T+1)
	for t := 1; t <= T; t++ {
		m := len(tasks[t])
		prefix[t] = make([]int, m+1)
		for j := 0; j < m; j++ {
			prefix[t][j+1] = prefix[t][j] + tasks[t][j]
		}
	}

	// dp[avail] = max interest achievable with `avail` free leaf slots remaining.
	// avail is capped at cap (meaning "more than enough slots").
	const cap = 1001
	const neg = -1
	dp := make([]int, cap+1)
	for i := range dp {
		dp[i] = neg
	}
	dp[1] = 0 // start: 1 slot at depth 0

	// Process from t=T (depth 0) down to t=1 (depth T-1).
	// At each level, choose k tasks (best k by interest); remaining slots double.
	for t := T; t >= 1; t-- {
		ndp := make([]int, cap+1)
		for i := range ndp {
			ndp[i] = neg
		}
		for avail := 0; avail <= cap; avail++ {
			if dp[avail] == neg {
				continue
			}
			maxK := len(tasks[t])
			if avail < maxK {
				maxK = avail
			}
			for k := 0; k <= maxK; k++ {
				newAvail := (avail - k) * 2
				if newAvail > cap {
					newAvail = cap
				}
				val := dp[avail] + prefix[t][k]
				if ndp[newAvail] < val {
					ndp[newAvail] = val
				}
			}
		}
		dp = ndp
	}

	ans := 0
	for avail := 0; avail <= cap; avail++ {
		if dp[avail] > ans {
			ans = dp[avail]
		}
	}
	fmt.Fprintln(out, ans)
}
