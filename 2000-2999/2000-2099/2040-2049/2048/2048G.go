package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const maxN = 1000000

var fact = make([]int64, maxN+5)
var invFact = make([]int64, maxN+5)

func modPow(base, exp int64) int64 {
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func prepare() {
	fact[0] = 1
	for i := 1; i < len(fact); i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[len(fact)-1] = modPow(fact[len(fact)-1], mod-2)
	for i := len(fact) - 1; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	prepare()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, v int
		var m int64
		fmt.Fscan(in, &n, &m, &v)

		powV := make([]int64, n+1)
		powV[0] = 1
		for i := 1; i <= n; i++ {
			powV[i] = powV[i-1] * int64(v) % mod
		}

		total := modPow(int64(v), int64(n)*m)
		badTotal := int64(0)

		powA := make([]int64, n+1)
		powC1 := make([]int64, n+1)
		powC2 := make([]int64, n+1)

		for t := 2; t <= v; t++ {
			a := int64(t - 1)
			base1 := int64(v - t + 1)
			base2 := int64(v - t + 2)

			powA[0] = 1
			for i := 1; i <= n; i++ {
				powA[i] = powA[i-1] * a % mod
			}

			powC1[0] = 1
			powC2[0] = 1
			for i := 1; i <= n; i++ {
				powC1[i] = powC1[i-1] * base1 % mod
				powC2[i] = powC2[i-1] * base2 % mod
			}

			// count for columns min <= t-1
			count1 := modPow((powV[n]-powC1[n]+mod)%mod, m)
			for r := 1; r <= n; r++ {
				base := powA[r] * powV[n-r] % mod
				term := modPow(base, m)
				coef := comb(n, r)
				if r%2 == 1 {
					count1 = (count1 - coef*term) % mod
				} else {
					count1 = (count1 + coef*term) % mod
				}
			}
			if count1 < 0 {
				count1 += mod
			}

			// count for columns min <= t-2
			count2 := int64(0)
			for r := 0; r <= n; r++ {
				base := (powA[r]*powV[n-r] - powC2[n-r]) % mod
				if base < 0 {
					base += mod
				}
				term := modPow(base, m)
				coef := comb(n, r)
				if r%2 == 1 {
					count2 = (count2 - coef*term) % mod
				} else {
					count2 = (count2 + coef*term) % mod
				}
			}
			if count2 < 0 {
				count2 += mod
			}

			bad := (count1 - count2) % mod
			if bad < 0 {
				bad += mod
			}
			badTotal = (badTotal + bad) % mod
		}

		ans := (total - badTotal) % mod
		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(out, ans)
	}
}
