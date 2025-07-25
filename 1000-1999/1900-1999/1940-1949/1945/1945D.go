package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		solve()
	}
}

func solve() {
	var n, m int
	fmt.Fscan(in, &n, &m)
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + b[i]
	}

	dp := make([]int64, n+2)
	best := pref[n] // dp[n+1] = 0
	for j := n; j >= 1; j-- {
		dp[j] = a[j] - pref[j] + best
		if v := dp[j] + pref[j-1]; v < best {
			best = v
		}
	}
	ans := dp[1]
	for i := 1; i <= m; i++ {
		if dp[i] < ans {
			ans = dp[i]
		}
	}
	fmt.Fprintln(out, ans)
}
