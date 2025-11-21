package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxN = 1000000
	mod  = 998244353
)

var spf [maxN + 1]int

func init() {
	primes := make([]int, 0)
	for i := 2; i <= maxN; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > maxN {
				break
			}
			spf[i*p] = p
		}
	}
}

func factorize(n int) map[int]int {
	res := make(map[int]int)
	if n <= 1 {
		return res
	}
	for n > 1 {
		p := spf[n]
		if p == 0 {
			p = n
		}
		cnt := 0
		for n%p == 0 {
			n /= p
			cnt++
		}
		res[p] = cnt
	}
	return res
}

func modPow(a int64, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow((a%mod+mod)%mod, mod-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, z int64
		fmt.Fscan(in, &x, &y, &z)

		mFactors := make(map[int]int)
		for _, v := range []int64{x, y, z} {
			fac := factorize(int(v))
			for p, e := range fac {
				mFactors[p] += e
			}
		}

		expMap := make(map[int]int)
		for p := range mFactors {
			pf := factorize(p - 1)
			for prime, exp := range pf {
				if _, exists := mFactors[prime]; exists {
					continue
				}
				expMap[prime] += exp
			}
		}

		total := int64(1)
		for prime, exp := range expMap {
			term := (modPow(int64(prime), int64(exp)) - 1 + mod) % mod
			term = term * modInv(int64(prime)) % mod
			term = (term + 1) % mod
			total = total * term % mod
		}

		mMod := (x % mod) * (y % mod) % mod
		mMod = mMod * (z % mod) % mod
		ans := total * modInv(mMod) % mod
		fmt.Fprintln(out, ans)
	}
}
