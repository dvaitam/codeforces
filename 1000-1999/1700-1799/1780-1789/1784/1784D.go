package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353
const ROOT int = 3

func modAdd(a, b int) int {
	s := a + b
	if s >= MOD {
		s -= MOD
	}
	return s
}
func modSub(a, b int) int {
	s := a - b
	if s < 0 {
		s += MOD
	}
	return s
}
func modMul(a, b int) int {
	return int(int64(a) * int64(b) % int64(MOD))
}
func modPow(a, e int) int {
	res := 1
	x := a
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, x)
		}
		x = modMul(x, x)
		e >>= 1
	}
	return res
}
func modInv(a int) int {
	return modPow(a, MOD-2)
}

func ntt(a []int, invert bool) {
	n := len(a)
	j := 0
	for i := 1; i < n; i++ {
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
		for i := 0; i < n; i++ {
			a[i] = modMul(a[i], invN)
		}
	}
}

func polyMul(a, b []int) []int {
	n := len(a) + len(b) - 1
	sz := 1
	for sz < n {
		sz <<= 1
	}
	fa := make([]int, sz)
	fb := make([]int, sz)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < sz; i++ {
		fa[i] = modMul(fa[i], fb[i])
	}
	ntt(fa, true)
	return fa[:n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m := 1 << uint(n)
	fact := make([]int, m+1)
	invfact := make([]int, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = modMul(fact[i-1], i)
	}
	invfact[m] = modInv(fact[m])
	for i := m; i > 0; i-- {
		invfact[i-1] = modMul(invfact[i], i)
	}
	if n == 1 {
		fmt.Println("0 2")
		return
	}
	Fprev := make([]int, 1<<1+1)
	Fprev[2] = 2
	for level := 2; level <= n; level++ {
		mcur := 1 << uint(level)
		half := mcur >> 1
		A := make([]int, half+1)
		for r := 1; r <= half; r++ {
			if r < len(Fprev) {
				val := Fprev[r]
				if val != 0 {
					A[r] = modMul(modMul(val, invfact[r-1]), invfact[half-r])
				}
			}
		}
		B := make([]int, half)
		for k := 0; k < half; k++ {
			B[k] = modMul(invfact[k], invfact[half-1-k])
		}
		C := polyMul(A, B)
		Fcurr := make([]int, mcur+1)
		factorPre := modMul(mcur%MOD, fact[half-1])
		for j := 2; j <= mcur; j++ {
			if j-1 >= len(C) {
				continue
			}
			val := C[j-1]
			if val == 0 {
				continue
			}
			tmp := modMul(fact[j-2], fact[mcur-j])
			tmp = modMul(tmp, val)
			tmp = modMul(tmp, factorPre)
			Fcurr[j] = tmp
		}
		Fprev = Fcurr
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= m; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, Fprev[i]%MOD)
	}
	fmt.Fprintln(out)
}
