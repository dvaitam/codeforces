package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

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
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] > maxA {
			maxA = arr[i]
		}
	}

	// factorials for combinations
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	comb := func(nn, kk int) int64 {
		if kk < 0 || kk > nn {
			return 0
		}
		return fact[nn] * invFact[kk] % MOD * invFact[nn-kk] % MOD
	}

	// powers of 2
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	// prefix sums of C(n-1, k)
	prefixC := make([]int64, n)
	for k := 0; k < n; k++ {
		if k == 0 {
			prefixC[k] = comb(n-1, 0)
		} else {
			prefixC[k] = (prefixC[k-1] + comb(n-1, k)) % MOD
		}
	}

	// compute D[i]
	D := make([]int64, n)
	for i := 0; i < n; i++ {
		R := n - 1 - i
		pre := int64(0)
		if R-1 >= 0 {
			pre = prefixC[R-1]
		}
		val := (pow2[n-1] - comb(n-1, R) - 2*pre) % MOD
		if val < 0 {
			val += MOD
		}
		D[i] = val
	}
	prefixD := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixD[i+1] = (prefixD[i] + D[i]) % MOD
	}

	// sieve for smallest prime factor
	spf := make([]int, maxA+1)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	// map prime -> slice of counts per exponent
	counts := make(map[int][]int)
	for _, x := range arr {
		v := x
		for v > 1 {
			p := spf[v]
			e := 0
			for v%p == 0 {
				v /= p
				e++
			}
			c := counts[p]
			if len(c) <= e {
				nc := make([]int, e+1)
				copy(nc, c)
				c = nc
			}
			c[e]++
			counts[p] = c
		}
	}

	var ans int64
	for _, c := range counts {
		totalPos := 0
		for _, v := range c {
			totalPos += v
		}
		count0 := n - totalPos
		pos := count0
		for e := 1; e < len(c); e++ {
			cnt := c[e]
			if cnt == 0 {
				continue
			}
			sumRange := (prefixD[pos+cnt] - prefixD[pos]) % MOD
			if sumRange < 0 {
				sumRange += MOD
			}
			ans = (ans + int64(e)%MOD*sumRange) % MOD
			pos += cnt
		}
	}

	fmt.Println(ans % MOD)
}
