package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var m, v int64
	if _, err := fmt.Fscan(in, &n, &m, &v); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		a[i] %= MOD
	}

	dp := make([]int64, n+1)
	dp[0] = 1
	for j := 1; j <= n; j++ {
		ndp := make([]int64, n+1)
		aj := a[j-1]
		for s := 0; s < j; s++ {
			val := dp[s]
			if val == 0 {
				continue
			}
			// j not selected
			ndp[s] = (ndp[s] + val*aj) % MOD
			// start new block with j
			ndp[s+1] = (ndp[s+1] + val*v%MOD*int64(j)) % MOD
			// join one of existing blocks
			if s > 0 {
				ndp[s] = (ndp[s] + val*v%MOD*int64(s)) % MOD
			}
		}
		dp = ndp
	}

	fall := make([]int64, n+1)
	fall[0] = 1
	for i := 1; i <= n; i++ {
		if m-int64(i)+1 <= 0 {
			fall[i] = 0
		} else {
			fall[i] = fall[i-1] * ((m - int64(i) + 1) % MOD) % MOD
		}
	}

	powN := make([]int64, n+1)
	powN[0] = 1
	for i := 1; i <= n; i++ {
		powN[i] = powN[i-1] * int64(n) % MOD
	}

	ans := int64(0)
	for s := 0; s <= n; s++ {
		if dp[s] == 0 {
			continue
		}
		ff := fall[s]
		inv := modPow(powN[s], MOD-2)
		ans = (ans + dp[s]%MOD*ff%MOD*inv) % MOD
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
