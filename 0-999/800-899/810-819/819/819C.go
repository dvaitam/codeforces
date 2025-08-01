package main

import (
	"bufio"
	"fmt"
	"os"
)

var primes []int64

func initPrimes() {
	const maxP = 1000000
	sieve := make([]bool, maxP+1)
	for i := 2; i*i <= maxP; i++ {
		if !sieve[i] {
			for j := i * i; j <= maxP; j += i {
				sieve[j] = true
			}
		}
	}
	for i := 2; i <= maxP; i++ {
		if !sieve[i] {
			primes = append(primes, int64(i))
		}
	}
}

func factor(x int64) map[int64]int {
	res := make(map[int64]int)
	for _, p := range primes {
		if p*p > x {
			break
		}
		for x%p == 0 {
			res[p]++
			x /= p
		}
	}
	if x > 1 {
		res[x]++
	}
	return res
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func powInt(x int64, e int) int64 {
	res := int64(1)
	for e > 0 {
		res *= x
		e--
	}
	return res
}

func coprimeCount(M int64, primes []int64) int64 {
	var dfs func(int, int64, int)
	ans := int64(0)
	dfs = func(i int, prod int64, sign int) {
		if i == len(primes) {
			if sign == 0 {
				ans += M / prod
			} else {
				ans += int64(sign) * (M / prod)
			}
			return
		}
		// skip current prime
		dfs(i+1, prod, sign)
		// take current prime
		if prod <= M/primes[i] {
			if sign == 0 {
				dfs(i+1, prod*primes[i], 1)
			} else {
				dfs(i+1, prod*primes[i], -sign)
			}
		}
	}
	dfs(0, 1, 0)
	return ans
}

func main() {
	initPrimes()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n1, n2, n3 int64
		fmt.Fscan(in, &n1, &n2, &n3)
		var m1, m2, m3 int64
		fmt.Fscan(in, &m1, &m2, &m3)
		var s1, s2, s3 int64
		fmt.Fscan(in, &s1, &s2, &s3)

		n := n1 * n2 * n3
		m := m1 * m2 * m3
		s := s1 * s2 * s3

		// factorization of n
		nf := factor(n1)
		for k, v := range factor(n2) {
			nf[k] += v
		}
		for k, v := range factor(n3) {
			nf[k] += v
		}
		// factorization of 2*s
		sf := factor(s1)
		for k, v := range factor(s2) {
			sf[k] += v
		}
		for k, v := range factor(s3) {
			sf[k] += v
		}
		sf[2]++ // multiply by 2
		twoS := s * 2

		g := gcd(n, twoS)
		// prepare prime-exponent arrays for n
		primesN := make([]int64, 0, len(nf))
		expsN := make([]int, 0, len(nf))
		index := make(map[int64]int)
		i := 0
		for p, e := range nf {
			primesN = append(primesN, p)
			expsN = append(expsN, e)
			index[p] = i
			i++
		}
		// g factors
		gf := make(map[int64]int)
		for p, e := range nf {
			if se, ok := sf[p]; ok {
				if e < se {
					gf[p] = e
				} else {
					gf[p] = se
				}
			}
		}
		primesG := make([]int64, 0, len(gf))
		expsG := make([]int, 0, len(gf))
		for p, e := range gf {
			primesG = append(primesG, p)
			expsG = append(expsG, e)
		}
		expsD := make([]int, len(primesN))
		var ans int64
		if g == n {
			ans++ // x=0 case
		}
		var dfs func(int, int64)
		dfs = func(pos int, curr int64) {
			if pos == len(primesG) {
				if curr > m {
					return
				}
				M := m / curr
				if M == 0 {
					return
				}
				var rest []int64
				for idx, p := range primesN {
					cnt := expsN[idx] - expsD[idx]
					if cnt > 0 {
						rest = append(rest, p)
					}
				}
				ans += coprimeCount(M, rest)
				return
			}
			p := primesG[pos]
			idx := index[p]
			power := int64(1)
			for k := 0; k <= expsG[pos]; k++ {
				expsD[idx] = k
				dfs(pos+1, curr*power)
				power *= p
			}
			expsD[idx] = 0
		}
		dfs(0, 1)
		// Stage2: divisors of twoS <= n-1
		// factorization of 2s is sf
		primesS := make([]int64, 0, len(sf))
		expsS := make([]int, 0, len(sf))
		for p, e := range sf {
			primesS = append(primesS, p)
			expsS = append(expsS, e)
		}
		var countDiv func(int, int64) int64
		limit := n - 1
		countDiv = func(pos int, prod int64) int64 {
			if prod > limit {
				return 0
			}
			if pos == len(primesS) {
				if prod >= 1 {
					return 1
				}
				return 0
			}
			res := int64(0)
			p := primesS[pos]
			val := int64(1)
			for e := 0; e <= expsS[pos]; e++ {
				res += countDiv(pos+1, prod*val)
				val *= p
				if prod*val > limit {
					break
				}
			}
			return res
		}
		if limit > 0 {
			ans += countDiv(0, 1)
		}
		fmt.Fprintln(out, ans)
	}
}
