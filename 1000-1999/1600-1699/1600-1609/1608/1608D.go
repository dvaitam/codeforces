package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func modPow(x, y int) int {
	res := 1
	for y > 0 {
		if y&1 == 1 {
			res = res * x % MOD
		}
		x = x * x % MOD
		y >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type dom struct{ l, r byte }
	ds := make([]dom, n)

	a, d := 0, 0 // counts of fixed BB and WW
	nb, nw, u := 0, 0, 0
	countNoCross := 1
	canBW, canWB := true, true

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		l, r := s[0], s[1]
		ds[i] = dom{l, r}
		switch {
		case l == 'B' && r == 'B':
			a++
		case l == 'W' && r == 'W':
			d++
		case l == 'B' && r == '?':
			nb++
		case l == '?' && r == 'B':
			nb++
		case l == 'W' && r == '?':
			nw++
		case l == '?' && r == 'W':
			nw++
		case l == '?' && r == '?':
			u++
		}
		options := 0
		if (l == 'B' || l == '?') && (r == 'W' || r == '?') {
			options++
		}
		if (l == 'W' || l == '?') && (r == 'B' || r == '?') {
			options++
		}
		if options == 0 {
			countNoCross = 0
		} else {
			countNoCross = countNoCross * options % MOD
		}
		if !(l == 'B' || l == '?') || !(r == 'W' || r == '?') {
			canBW = false
		}
		if !(l == 'W' || l == '?') || !(r == 'B' || r == '?') {
			canWB = false
		}
	}

	maxN := 2*n + 5
	fact := make([]int, maxN)
	inv := make([]int, maxN)
	fact[0] = 1
	for i := 1; i < maxN; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	inv[maxN-1] = modPow(fact[maxN-1], MOD-2)
	for i := maxN - 2; i >= 0; i-- {
		inv[i] = inv[i+1] * (i + 1) % MOD
	}
	comb := func(n, k int) int {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * inv[k] % MOD * inv[n-k] % MOD
	}

	N := nb + nw + 2*u
	t := nw + u - (a - d)
	countEq := 0
	if t >= 0 && t <= N {
		countEq = comb(N, t)
	}

	ans := (countEq - countNoCross) % MOD
	if ans < 0 {
		ans += MOD
	}
	if canBW {
		ans = (ans + 1) % MOD
	}
	if canWB {
		ans = (ans + 1) % MOD
	}
	fmt.Println(ans)
}
