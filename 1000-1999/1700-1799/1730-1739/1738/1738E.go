package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXN = 100000

var fact [MAXN + 1]int64
var invFact [MAXN + 1]int64

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
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
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func solve(n int, a []int64) int64 {
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}
	total := prefix[n]
	if total == 0 {
		return modPow(2, int64(n-1))
	}

	cnt := make(map[int64]int)
	for i := 1; i < n; i++ {
		cnt[prefix[i]]++
	}

	ans := int64(1)
	for s, l := range cnt {
		if l == 0 || 2*s >= total {
			continue
		}
		r := cnt[total-s]
		ans = ans * comb(int64(l+r), int64(l)) % MOD
		cnt[s] = 0
		cnt[total-s] = 0
	}
	if total%2 == 0 {
		c := cnt[total/2]
		if c > 0 {
			ans = ans * modPow(2, int64(c)) % MOD
		}
	}
	return ans % MOD
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
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		fmt.Fprintln(out, solve(n, arr))
	}
}
