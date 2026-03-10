package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func solve(n int) int64 {
	if n == 2 {
		return 1
	}
	dp := make([]int64, n+1)
	s := int64(0)
	for i := n; i >= 1; i-- {
		dp[i] = (int64(i)%mod*s%mod + 1) % mod
		s = (s + dp[i]) % mod
	}
	ans := int64(n - 1)
	for k := 3; k <= n; k++ {
		ways := int64((k-1)*(k-2)/2 - 1)
		ans = (ans + ways%mod*dp[k]%mod) % mod
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, solve(n))
	}
}
