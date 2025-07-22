package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN = 100000

var mu [MAXN + 1]int
var fac, invfac [MAXN + 1]int64
var divisors [MAXN + 1][]int

func modexp(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		e >>= 1
	}
	return res
}

func initMuAndDivs() {
	mu[1] = 1
	primes := make([]int, 0, MAXN/10)
	isComp := make([]bool, MAXN+1)
	for i := 2; i <= MAXN; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			ip := i * p
			if ip > MAXN {
				break
			}
			isComp[ip] = true
			if i%p == 0 {
				mu[ip] = 0
				break
			} else {
				mu[ip] = -mu[i]
			}
		}
	}
	for i := 1; i <= MAXN; i++ {
		for j := i; j <= MAXN; j += i {
			divisors[j] = append(divisors[j], i)
		}
	}
}

func initFactorials() {
	fac[0] = 1
	for i := 1; i <= MAXN; i++ {
		fac[i] = (fac[i-1] * int64(i)) % MOD
	}
	invfac[MAXN] = modexp(fac[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		invfac[i-1] = (invfac[i] * int64(i)) % MOD
	}
}

// comb returns C(n, k)
func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * invfac[k] % MOD * invfac[n-k] % MOD
}

func main() {
	initMuAndDivs()
	initFactorials()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var n, f int
		fmt.Fscan(reader, &n, &f)
		var ans int64
		// sum over d|n mu[d] * C(n/d -1, f-1)
		for _, d := range divisors[n] {
			nd := n / d
			if nd < f {
				continue
			}
			c := comb(nd-1, f-1)
			m := mu[d]
			if m == 0 || c == 0 {
				continue
			}
			ans = (ans + int64(m)*c%MOD + MOD) % MOD
		}
		fmt.Fprintln(writer, ans)
	}
}
