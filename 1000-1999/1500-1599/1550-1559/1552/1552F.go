package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	x := make([]int64, n)
	y := make([]int64, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &y[i], &s[i])
	}

	dp := make([]int64, n)
	pref := make([]int64, n+1)
	ans := (x[n-1] + 1) % mod
	for i := 0; i < n; i++ {
		idx := sort.Search(len(x), func(j int) bool { return x[j] >= y[i] })
		val := (x[i] - y[i] + pref[i] - pref[idx]) % mod
		if val < 0 {
			val += mod
		}
		dp[i] = val
		pref[i+1] = (pref[i] + val) % mod
		if s[i] == 1 {
			ans = (ans + val) % mod
		}
	}
	fmt.Println(ans % mod)
}
