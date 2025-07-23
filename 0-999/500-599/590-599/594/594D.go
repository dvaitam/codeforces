package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n, tree: make([]int64, n+1)}
	for i := 0; i <= n; i++ {
		b.tree[i] = 1
	}
	return b
}

func (b *BIT) Mul(i int, v int64) {
	for i <= b.n {
		b.tree[i] = b.tree[i] * v % mod
		i += i & -i
	}
}

func (b *BIT) Pref(i int) int64 {
	res := int64(1)
	for i > 0 {
		res = res * b.tree[i] % mod
		i -= i & -i
	}
	return res
}

func powmod(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 != 0 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func inv(a int64) int64 {
	return powmod(a, mod-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	maxA := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}

	var q int
	fmt.Fscan(in, &q)
	type Query struct{ l, idx int }
	queries := make([][]Query, n+1)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		queries[r] = append(queries[r], Query{l: l, idx: i})
	}

	// prefix products
	pref := make([]int64, n+1)
	pref[0] = 1
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] * int64(a[i]) % mod
	}

	// compute smallest prime factors up to maxA
	spf := make([]int, maxA+1)
	primes := make([]int, 0)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > maxA {
				break
			}
			spf[i*p] = p
		}
	}

	// factorization helper
	uniqPrimes := func(x int) []int {
		res := []int{}
		for x > 1 {
			p := spf[x]
			res = append(res, p)
			for x%p == 0 {
				x /= p
			}
		}
		return res
	}

	bit := NewBIT(n)
	last := make([]int, maxA+1)
	ans := make([]int64, q)

	for i := 1; i <= n; i++ {
		for _, p := range uniqPrimes(a[i]) {
			f := (int64(p-1) * inv(int64(p))) % mod
			if last[p] != 0 {
				bit.Mul(last[p], inv(f))
			}
			bit.Mul(i, f)
			last[p] = i
		}

		for _, qu := range queries[i] {
			l := qu.l
			prod := pref[i] * inv(pref[l-1]) % mod
			coef := bit.Pref(i) * inv(bit.Pref(l-1)) % mod
			ans[qu.idx] = prod * coef % mod
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
