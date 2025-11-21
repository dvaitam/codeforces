package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const primitiveRoot int64 = 3

func modPow(a, e int64) int64 {
	a %= mod
	if a < 0 {
		a += mod
	}
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
		step := (mod - 1) / int64(length)
		wlen := modPow(primitiveRoot, step)
		if invert {
			wlen = modPow(wlen, mod-2)
		}
		for i := 0; i < n; i += length {
			w := int64(1)
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := a[i+j+half] * w % mod
				a[i+j] = (u + v) % mod
				a[i+j+half] = (u - v + mod) % mod
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

func convolution(a, b []int64, limit int) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return []int64{}
	}
	if len(a) > limit {
		a = a[:limit]
	}
	if len(b) > limit {
		b = b[:limit]
	}
	size := 1
	need := len(a) + len(b) - 1
	for size < need {
		size <<= 1
	}
	fa := make([]int64, size)
	fb := make([]int64, size)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < size; i++ {
		fa[i] = fa[i] * fb[i] % mod
	}
	ntt(fa, true)
	if need > limit {
		need = limit
	}
	res := make([]int64, need)
	copy(res, fa[:need])
	return res
}

func buildSeries(cnt int, alpha int64, D int, inv []int64) []int64 {
	res := make([]int64, D+1)
	coeff := int64(1)
	alphaPow := int64(1)
	for x := 0; x <= D; x++ {
		if x == 0 {
			coeff = 1
			alphaPow = 1
		} else {
			coeff = coeff * int64(x+cnt-1) % mod * inv[x] % mod
			alphaPow = alphaPow * alpha % mod
		}
		res[x] = coeff * alphaPow % mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	lens := make([]int, k)
	cntGeom := make([]int, 6)
	totalLen := 0
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &lens[i])
		totalLen += lens[i]
		for t := 1; t <= lens[i]; t++ {
			cntGeom[t]++
		}
	}
	// consume arrays, values are irrelevant for the count
	var dump int
	for i := 0; i < k; i++ {
		for j := 0; j < lens[i]; j++ {
			fmt.Fscan(in, &dump)
		}
	}
	if totalLen > n {
		fmt.Println(0)
		return
	}
	D := n - totalLen
	inv := make([]int64, D+2)
	if D >= 1 {
		inv[1] = 1
		for i := 2; i <= D; i++ {
			idx := int(mod % int64(i))
			inv[i] = (mod - (mod/int64(i))*inv[idx]%mod) % mod
		}
	}
	invM := modPow(int64(m%int(mod)), mod-2)
	pol := []int64{1}
	limit := D + 1
	for t := 1; t <= 5; t++ {
		c := cntGeom[t]
		if c == 0 {
			continue
		}
		alpha := int64(m - t)
		alpha %= mod
		if alpha < 0 {
			alpha += mod
		}
		alpha = alpha * invM % mod
		series := buildSeries(c, alpha, D, inv)
		pol = convolution(pol, series, limit)
	}
	sumP := int64(0)
	for _, v := range pol {
		sumP += v
		if sumP >= mod {
			sumP -= mod
		}
	}
	powM := modPow(int64(m%int(mod)), int64(D))
	prod := int64(1)
	for t := 1; t <= 5; t++ {
		if cntGeom[t] == 0 {
			continue
		}
		prod = prod * modPow(int64(t), int64(cntGeom[t])) % mod
	}
	ans := powM * prod % mod * sumP % mod
	fmt.Println(ans)
}
