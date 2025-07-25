package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func powMod(x, e int64) int64 {
	res := int64(1)
	x %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * x % MOD
		}
		x = x * x % MOD
		e >>= 1
	}
	return res
}

func fwtInt(a []int64) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				u := a[i+j]
				v := a[i+j+step]
				a[i+j] = u + v
				a[i+j+step] = u - v
			}
		}
	}
}

func fwt(a []int64, inv bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				u := a[i+j]
				v := a[i+j+step]
				a[i+j] = (u + v) % MOD
				a[i+j+step] = (u - v) % MOD
				if a[i+j+step] < 0 {
					a[i+j+step] += MOD
				}
			}
		}
	}
	if inv {
		invN := powMod(int64(n), MOD-2)
		for i := 0; i < n; i++ {
			a[i] = a[i] * invN % MOD
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	vals := make([]int64, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		vals[i] = int64(x)
		if x > maxVal {
			maxVal = x
		}
	}

	size := 1
	for size <= maxVal {
		size <<= 1
	}
	if size < 1<<17 {
		size = 1 << 17
	}

	freq := make([]int64, size)
	for _, v := range vals {
		freq[v]++
	}

	fwtInt(freq)

	res := make([]int64, size)
	N := int64(n)
	for i := 0; i < size; i++ {
		val := freq[i]
		even := (N + val) / 2
		odd := N - even
		res[i] = powMod(3, even) * powMod(MOD-1, odd) % MOD
	}

	fwt(res, true)

	fmt.Fprintln(out, res[0]%MOD)
}
