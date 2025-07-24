package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func modAdd(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += mod
	}
	return a
}
func modMul(a, b int) int { return int((int64(a) * int64(b)) % int64(mod)) }

func modPow(a, e int) int {
	res := 1
	x := a % mod
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, x)
		}
		x = modMul(x, x)
		e >>= 1
	}
	return res
}

func modInv(a int) int { return modPow(a, mod-2) }

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
		wlen := modPow(3, (mod-1)/length)
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

func multiply(a, b []int) []int {
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

func multiplyTrunc(a, b []int, limit int) []int {
	res := multiply(a, b)
	if len(res) > limit {
		return res[:limit]
	}
	return res
}

var K int

func buildPoly(vals []int, l, r int) []int {
	if l == r {
		return []int{vals[l], 1}
	}
	m := (l + r) / 2
	left := buildPoly(vals, l, m)
	right := buildPoly(vals, m+1, r)
	return multiplyTrunc(left, right, K+1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n, &K)
	base := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &base[i])
	}
	var qnum int
	fmt.Fscan(in, &qnum)
	out := bufio.NewWriter(os.Stdout)
	for ; qnum > 0; qnum-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var q int64
			var idx int
			var d int64
			fmt.Fscan(in, &q, &idx, &d)
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				val := base[i]
				if i == idx-1 {
					val = d
				}
				x := ((q-val)%int64(mod) + int64(mod)) % int64(mod)
				arr[i] = int(x)
			}
			poly := buildPoly(arr, 0, n-1)
			if K < len(poly) {
				fmt.Fprintln(out, poly[K]%mod)
			} else {
				fmt.Fprintln(out, 0)
			}
		} else {
			var q int64
			var L, R int
			var d int64
			fmt.Fscan(in, &q, &L, &R, &d)
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				val := base[i]
				if i+1 >= L && i+1 <= R {
					val += d
				}
				x := ((q-val)%int64(mod) + int64(mod)) % int64(mod)
				arr[i] = int(x)
			}
			poly := buildPoly(arr, 0, n-1)
			if K < len(poly) {
				fmt.Fprintln(out, poly[K]%mod)
			} else {
				fmt.Fprintln(out, 0)
			}
		}
	}
	out.Flush()
}
