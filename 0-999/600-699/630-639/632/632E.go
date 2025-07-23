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
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, mod-2)
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
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := int64(1)
			for j := 0; j < length/2; j++ {
				u := a[i+j]
				v := a[i+j+length/2] * w % mod
				a[i+j] = (u + v) % mod
				a[i+j+length/2] = (u - v + mod) % mod
				w = w * wlen % mod
			}
		}
	}
	if invert {
		invN := modInv(int64(n))
		for i := 0; i < n; i++ {
			a[i] = a[i] * invN % mod
		}
	}
}

func convNTT(a, b []int64, limit int) []int64 {
	n := 1
	for n < len(a)+len(b) {
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
	resLen := len(a) + len(b) - 1
	if resLen > limit+1 {
		resLen = limit + 1
	}
	res := make([]int64, resLen)
	for i := 0; i < resLen; i++ {
		if fa[i]%mod != 0 {
			res[i] = 1
		}
	}
	return res
}

func powPoly(base []int64, k, limit int) []int64 {
	res := make([]int64, 1)
	res[0] = 1
	for k > 0 {
		if k&1 == 1 {
			res = convNTT(res, base, limit)
		}
		k >>= 1
		if k > 0 {
			base = convNTT(base, base, limit)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	costs := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &costs[i])
		if costs[i] > maxA {
			maxA = costs[i]
		}
	}
	limit := k * maxA
	base := make([]int64, limit+1)
	for _, v := range costs {
		if v <= limit {
			base[v] = 1
		}
	}
	res := powPoly(base, k, limit)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	first := true
	for i := 0; i < len(res); i++ {
		if res[i] != 0 {
			if !first {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, i)
			first = false
		}
	}
	if first {
		fmt.Fprintln(out)
	} else {
		fmt.Fprintln(out)
	}
}
