package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN int = 1000020

var fac [MAXN + 1]int64
var ifac [MAXN + 1]int64
var primes []int

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

func initFactorials() {
	fac[0] = 1
	for i := 1; i <= MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[MAXN] = modPow(fac[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func sievePrimes(limit int) {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= limit; p++ {
		if isPrime[p] {
			for j := p * p; j <= limit; j += p {
				isPrime[j] = false
			}
		}
	}
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
}

type pair struct {
	p int
	e int
}

func factorize(x int) []pair {
	res := []pair{}
	tmp := x
	for _, p := range primes {
		if p*p > tmp {
			break
		}
		if tmp%p == 0 {
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			res = append(res, pair{p, cnt})
		}
	}
	if tmp > 1 {
		res = append(res, pair{tmp, 1})
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	sievePrimes(1000000)
	initFactorials()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}

	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		factors := factorize(x)
		ans := modPow(2, int64(y-1))
		for _, pe := range factors {
			ans = ans * C(pe.e+y-1, y-1) % MOD
		}
		fmt.Fprintln(out, ans)
	}
}
