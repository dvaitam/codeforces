package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(x, y int) int {
	res := 1
	x %= mod
	for y > 0 {
		if y&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		y >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a, b, c int
	if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
		return
	}

	maxN := a
	if b > maxN {
		maxN = b
	}
	if c > maxN {
		maxN = c
	}

	fact := make([]int, maxN+1)
	invFact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}

	comb := func(n, k int) int {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}

	F := func(x, y int) int {
		limit := x
		if y < limit {
			limit = y
		}
		res := 0
		for i := 0; i <= limit; i++ {
			add := comb(x, i)
			add = add * comb(y, i) % mod
			add = add * fact[i] % mod
			res += add
			if res >= mod {
				res -= mod
			}
		}
		return res
	}

	ans := F(a, b)
	ans = ans * F(a, c) % mod
	ans = ans * F(b, c) % mod

	fmt.Fprintln(writer, ans)
}
