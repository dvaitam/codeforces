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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve(in, out)
	}
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n, k int
	fmt.Fscan(in, &n, &k)

	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	cold := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		fmt.Fscan(in, &cold[i])
	}
	hot := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		fmt.Fscan(in, &hot[i])
	}

	const INF int64 = 1 << 60
	dp := make([]int64, k+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = cold[a[1]]
	last := a[1]

	for i := 2; i <= n; i++ {
		ndp := make([]int64, k+1)
		for j := range ndp {
			ndp[j] = INF
		}
		ai := a[i]
		for other := 0; other <= k; other++ {
			cur := dp[other]
			if cur == INF {
				continue
			}
			costSame := cold[ai]
			if ai == last {
				costSame = hot[ai]
			}
			if cur+costSame < ndp[other] {
				ndp[other] = cur + costSame
			}
			costSwitch := cold[ai]
			if ai == other {
				costSwitch = hot[ai]
			}
			if cur+costSwitch < ndp[last] {
				ndp[last] = cur + costSwitch
			}
		}
		dp = ndp
		last = ai
	}

	ans := dp[0]
	for i := 1; i <= k; i++ {
		if dp[i] < ans {
			ans = dp[i]
		}
	}
	fmt.Fprintln(out, ans)
}
