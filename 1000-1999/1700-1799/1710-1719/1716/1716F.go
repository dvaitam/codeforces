package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXK = 2000

var stirling [MAXK + 1][MAXK + 1]int64

func init() {
	stirling[0][0] = 1
	for n := 1; n <= MAXK; n++ {
		for k := 1; k <= n; k++ {
			stirling[n][k] = (stirling[n-1][k-1] + int64(k)*stirling[n-1][k]) % MOD
		}
	}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		a := (m + 1) / 2
		invm := modPow(m%MOD, MOD-2)
		p := a % MOD * invm % MOD
		mPow := modPow(m%MOD, n)
		maxI := k
		if n < maxI {
			maxI = n
		}
		fall := int64(1)
		powp := int64(1)
		ans := stirling[k][0] * fall % MOD * powp % MOD
		for i := int64(1); i <= maxI; i++ {
			fall = fall * ((n - i + 1) % MOD) % MOD
			powp = powp * p % MOD
			ans = (ans + stirling[k][i]*fall%MOD*powp) % MOD
		}
		ans = ans * mPow % MOD
		fmt.Fprintln(writer, ans)
	}
}
