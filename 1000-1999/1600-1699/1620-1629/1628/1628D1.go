package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7
const MAXN int = 2000

var fact [MAXN + 1]int64
var invfact [MAXN + 1]int64
var pow2 [MAXN + 1]int64
var invPow2 [MAXN + 1]int64

func powmod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initComb() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[MAXN] = powmod(fact[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invfact[i-1] = invfact[i] * int64(i) % MOD
	}
	pow2[0] = 1
	for i := 1; i <= MAXN; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	invPow2[MAXN] = powmod(pow2[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invPow2[i-1] = invPow2[i] * 2 % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

func main() {
	initComb()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		var k int64
		fmt.Fscan(in, &n, &m, &k)
		sum := int64(0)
		for x := 0; x < m; x++ {
			sum = (sum + int64(m-x)*C(n, x)) % MOD
		}
		ans := k % MOD * sum % MOD * invPow2[n-1] % MOD
		fmt.Fprintln(out, ans)
	}
}
