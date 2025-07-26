package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func solve(a []int) int64 {
	n := len(a)
	dp0 := make([]int64, n+5)
	dp1 := make([]int64, n+5)
	dp0[0] = 1
	for _, x := range a {
		dp0[x+1] = (dp0[x+1] * 2) % mod
		dp1[x+1] = (dp1[x+1] * 2) % mod

		tmp := dp0[x]
		dp0[x+1] = (dp0[x+1] + tmp) % mod

		if x >= 1 {
			tmp0 := dp0[x-1]
			tmp1 := dp1[x-1]
			dp1[x-1] = (tmp1*2 + tmp0) % mod
		}
	}
	var ans int64
	for i := 0; i < len(dp0); i++ {
		ans += dp0[i] + dp1[i]
	}
	ans = (ans - 1) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solve(a))
	}
}
