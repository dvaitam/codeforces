package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

var (
	phi  []int
	spf  []int
	fac  []int64
	ifac []int64
	cnt  []int
	k    int
	ans  int64
)

func modPow(a int64, e int64) int64 {
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

func comb(n, r int) int64 {
	if n < r || r < 0 {
		return 0
	}
	return fac[n] * ifac[r] % MOD * ifac[n-r] % MOD
}

type pair struct{ p, e int }

func getDivisors(x int) []int {
	factors := make([]pair, 0, 8)
	for x > 1 {
		p := spf[x]
		c := 0
		for x%p == 0 {
			x /= p
			c++
		}
		factors = append(factors, pair{p, c})
	}
	divs := []int{1}
	for _, f := range factors {
		size := len(divs)
		mul := 1
		for i := 1; i <= f.e; i++ {
			mul *= f.p
			for j := 0; j < size; j++ {
				divs = append(divs, divs[j]*mul)
			}
		}
	}
	return divs
}

func addValue(x int) {
	divs := getDivisors(x)
	for _, d := range divs {
		old := cnt[d]
		cnt[d]++
		if cnt[d] >= k { // only compute difference if cnt[d] >= k to save comb call
			delta := (comb(cnt[d], k) - comb(old, k)) % MOD
			if delta < 0 {
				delta += MOD
			}
			ans = (ans + int64(phi[d])*delta) % MOD
		}
	}
}

func sieve(maxn int) {
	phi = make([]int, maxn+1)
	spf = make([]int, maxn+1)
	phi[1] = 1
	primes := make([]int, 0)
	for i := 2; i <= maxn; i++ {
		if spf[i] == 0 {
			spf[i] = i
			phi[i] = i - 1
			primes = append(primes, i)
		}
		for _, p := range primes {
			v := p * i
			if v > maxn {
				break
			}
			spf[v] = p
			if i%p == 0 {
				phi[v] = phi[i] * p
				break
			} else {
				phi[v] = phi[i] * (p - 1)
			}
		}
	}
}

func prepareFact(n int) {
	fac = make([]int64, n+1)
	ifac = make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[n] = modPow(fac[n], MOD-2)
	for i := n; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &k, &q)
	values := make([]int, n+q)
	maxv := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
		if values[i] > maxv {
			maxv = values[i]
		}
	}
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &values[n+i])
		if values[n+i] > maxv {
			maxv = values[n+i]
		}
	}
	sieve(maxv)
	prepareFact(n + q)
	cnt = make([]int, maxv+1)
	// process initial values
	for i := 0; i < n; i++ {
		addValue(values[i])
	}
	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < q; i++ {
		addValue(values[n+i])
		fmt.Fprintln(out, ans%MOD)
	}
	out.Flush()
}
