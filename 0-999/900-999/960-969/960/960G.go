package main

import (
	"fmt"
)

const MOD int = 998244353
const ROOT int = 3

func modAdd(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}
func modMul(a, b int) int { return int(int64(a) * int64(b) % int64(MOD)) }
func modPow(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = modMul(res, base)
		}
		base = modMul(base, base)
		e >>= 1
	}
	return res
}
func modInv(a int) int { return modPow(a, MOD-2) }

func ntt(a []int, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j |= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(ROOT, (MOD-1)/length)
		if invert {
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := 1
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := modMul(a[i+j+half], w)
				a[i+j] = modAdd(u, v)
				a[i+j+half] = modSub(u, v)
				w = modMul(w, wlen)
			}
		}
	}
	if invert {
		invN := modInv(n)
		for i := range a {
			a[i] = modMul(a[i], invN)
		}
	}
}

func polyMul(a, b []int) []int {
	need := len(a) + len(b) - 1
	n := 1
	for n < need {
		n <<= 1
	}
	fa := make([]int, n)
	fb := make([]int, n)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < n; i++ {
		fa[i] = modMul(fa[i], fb[i])
	}
	ntt(fa, true)
	return fa[:need]
}

func build(l, r, k int) []int {
	if l == r {
		res := []int{l % MOD, 1}
		if len(res) > k+1 {
			res = res[:k+1]
		}
		return res
	}
	mid := (l + r) / 2
	left := build(l, mid, k)
	right := build(mid+1, r, k)
	res := polyMul(left, right)
	if len(res) > k+1 {
		res = res[:k+1]
	}
	return res
}

func stirlingFirst(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	poly := build(0, n-1, k)
	if k >= len(poly) {
		return 0
	}
	return poly[k]
}

func factPre(n int) ([]int, []int) {
	fac := make([]int, n+1)
	ifac := make([]int, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = modMul(fac[i-1], i)
	}
	ifac[n] = modInv(fac[n])
	for i := n; i > 0; i-- {
		ifac[i-1] = modMul(ifac[i], i)
	}
	return fac, ifac
}

func C(n, k int, fac, ifac []int) int {
	if k < 0 || k > n {
		return 0
	}
	return modMul(fac[n], modMul(ifac[k], ifac[n-k]))
}

func main() {
	var N, A, B int
	fmt.Scan(&N, &A, &B)
	if A == 0 || B == 0 || A+B-1 > N {
		fmt.Println(0)
		return
	}
	k := A + B - 2
	s := stirlingFirst(N-1, k)
	fac, ifac := factPre(k)
	comb := C(k, A-1, fac, ifac)
	ans := modMul(s, comb)
	fmt.Println(ans)
}
