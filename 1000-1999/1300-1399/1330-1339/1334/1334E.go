package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// modular exponentiation
func modPow(a, e, m int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var D int64
	if _, err := fmt.Fscan(reader, &D); err != nil {
		return
	}

	// factorize D
	primes := make([]int64, 0)
	exps := make([]int, 0)
	n := D
	for p := int64(2); p*p <= n; p++ {
		if n%p == 0 {
			cnt := 0
			for n%p == 0 {
				n /= p
				cnt++
			}
			primes = append(primes, p)
			exps = append(exps, cnt)
		}
	}
	if n > 1 {
		primes = append(primes, n)
		exps = append(exps, 1)
	}

	// determine factorial limit
	maxExp := 0
	for _, e := range exps {
		maxExp += e
	}
	// precompute factorials and inverses
	fact := make([]int64, maxExp+1)
	invFact := make([]int64, maxExp+1)
	fact[0] = 1
	for i := 1; i <= maxExp; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxExp] = modPow(fact[maxExp], mod-2, mod)
	for i := maxExp; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	var q int
	fmt.Fscan(reader, &q)

	for ; q > 0; q-- {
		var v, u int64
		fmt.Fscan(reader, &v, &u)
		// compute exponent counts relative to D's primes
		diffV := make([]int, len(primes))
		diffU := make([]int, len(primes))
		sumV := 0
		sumU := 0
		tempV := v
		tempU := u
		for i, p := range primes {
			cv := 0
			cu := 0
			for tempV%p == 0 {
				tempV /= p
				cv++
			}
			for tempU%p == 0 {
				tempU /= p
				cu++
			}
			g := cv
			if cu < g {
				g = cu
			}
			diffV[i] = cv - g
			diffU[i] = cu - g
			sumV += diffV[i]
			sumU += diffU[i]
		}

		// compute ways using factorials
		res := fact[sumV]
		for _, x := range diffV {
			res = res * invFact[x] % mod
		}
		res2 := fact[sumU]
		for _, x := range diffU {
			res2 = res2 * invFact[x] % mod
		}
		res = res * res2 % mod
		fmt.Fprintln(writer, res)
	}
}
