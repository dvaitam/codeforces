package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const MAXN = 1000000

var fact [MAXN + 1]int64
var invFact [MAXN + 1]int64
var phi [MAXN + 1]int

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

func init() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	for i := 0; i <= MAXN; i++ {
		phi[i] = i
	}
	for i := 2; i <= MAXN; i++ {
		if phi[i] == i {
			for j := i; j <= MAXN; j += i {
				phi[j] -= phi[j] / i
			}
		}
	}
}

func divisors(n int) []int {
	res := make([]int, 0)
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			res = append(res, i)
			if i*i != n {
				res = append(res, n/i)
			}
		}
	}
	sort.Ints(res)
	return res
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
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
		cnt := make([]int, n+1)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			cnt[v]++
		}
		// compute gcd of counts
		g := 0
		for i := 1; i <= n; i++ {
			if cnt[i] > 0 {
				if g == 0 {
					g = cnt[i]
				} else {
					g = gcd(g, cnt[i])
				}
			}
		}

		divs := divisors(n)
		var num, den int64
		for _, d := range divs { // d is divisor g
			m := n / d
			if g%m != 0 {
				continue
			}
			fix := fact[d]
			sumC := int64(0)
			for i := 1; i <= n; i++ {
				if cnt[i] > 0 {
					c := cnt[i] / m
					fix = fix * invFact[c] % MOD
					sumC += int64(c) * int64(c-1)
				}
			}
			coef := int64(phi[n/d])
			if d == 1 {
				// only one arrangement (all equal)
				den = (den + coef*fix) % MOD
				num = (num + coef*fix%MOD*1%MOD) % MOD
			} else {
				denomVal := int64(d) * int64(d-1) % MOD
				temp := sumC % MOD * modPow(denomVal, MOD-2) % MOD
				compSum := fix * int64(n) % MOD * ((1 - temp + MOD) % MOD) % MOD
				den = (den + coef*fix) % MOD
				num = (num + coef*compSum) % MOD
			}
		}
		ans := num % MOD * modPow(den%MOD, MOD-2) % MOD
		fmt.Fprintln(out, ans)
	}
}
