package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const root int64 = 3

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 != 0 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func ntt(a []int64, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(root, (mod-1)/int64(length))
		if invert {
			wlen = modPow(wlen, mod-2)
		}
		for i := 0; i < n; i += length {
			w := int64(1)
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := a[i+j+half] * w % mod
				a[i+j] = u + v
				if a[i+j] >= mod {
					a[i+j] -= mod
				}
				a[i+j+half] = u - v
				if a[i+j+half] < 0 {
					a[i+j+half] += mod
				}
				w = w * wlen % mod
			}
		}
	}
	if invert {
		invN := modPow(int64(n), mod-2)
		for i := range a {
			a[i] = a[i] * invN % mod
		}
	}
}

func convolution(a, b []int64) []int64 {
	n := 1
	need := len(a) + len(b) - 1
	for n < need {
		n <<= 1
	}
	fa := make([]int64, n)
	fb := make([]int64, n)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < n; i++ {
		fa[i] = fa[i] * fb[i] % mod
	}
	ntt(fa, true)
	return fa[:need]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	b := make([]int, n)
	maxN := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
		if a[i]+b[i] > maxN {
			maxN = a[i] + b[i]
		}
	}
	fac := make([]int64, maxN+1)
	ifac := make([]int64, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[maxN] = modPow(fac[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * ifac[k] % mod * ifac[n-k] % mod
	}
	dp := []int64{1}
	cur := 0
	for i := 0; i < n; i++ {
		ai, bi := a[i], b[i]
		newH := cur + ai - bi
		L := bi - cur
		if L < 0 {
			L = 0
		}
		R := bi + newH
		if R > ai+bi {
			R = ai + bi
		}
		arr2 := make([]int64, R-L+1)
		for k := L; k <= R; k++ {
			arr2[k-L] = comb(ai+bi, k)
		}
		arr1 := make([]int64, len(dp))
		copy(arr1, dp)
		conv := convolution(arr1, arr2)
		offset := bi - L
		dp = conv[offset : offset+newH+1]
		cur = newH
	}
	ans := int64(0)
	for _, v := range dp {
		ans += v
		if ans >= mod {
			ans -= mod
		}
	}
	fmt.Println(ans)
}
