package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

var fact []int64
var invFact []int64
var spf []int

func modPow(a, e int64) int64 {
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

func comb(n, k int) int64 {
	if n < k || k < 0 {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	const lim = 1000020
	fact = make([]int64, lim+1)
	invFact = make([]int64, lim+1)
	fact[0] = 1
	for i := 1; i <= lim; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[lim] = modPow(fact[lim], MOD-2)
	for i := lim; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	spf = make([]int, 1000000+1)
	for i := 2; i*i <= 1000000; i++ {
		if spf[i] == 0 {
			for j := i * i; j <= 1000000; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	for i := 2; i <= 1000000; i++ {
		if spf[i] == 0 {
			spf[i] = i
		}
	}

	for ; q > 0; q-- {
		var r, n int
		fmt.Fscan(in, &r, &n)
		ans := int64(1)
		m := n
		for m > 1 {
			p := spf[m]
			cnt := 0
			for m%p == 0 {
				m /= p
				cnt++
			}
			val := comb(r+cnt+1, r+1) - comb(r+cnt-1, r+1)
			if val < 0 {
				val += MOD
			}
			ans = ans * val % MOD
		}
		fmt.Fprintln(out, ans%MOD)
	}
}
