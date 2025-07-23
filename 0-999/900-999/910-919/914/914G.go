package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const BITS = 17
const SIZE = 1 << BITS

func addmod(a, b int64) int64 {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func submod(a, b int64) int64 {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}

func fwtXor(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				u := a[i+j]
				v := a[i+j+step]
				a[i+j] = addmod(u, v)
				a[i+j+step] = submod(u, v)
			}
		}
	}
	if invert {
		inv := powMod(int64(n), MOD-2)
		for i := 0; i < n; i++ {
			a[i] = a[i] * inv % MOD
		}
	}
}

func fwtOr(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				if !invert {
					a[i+j+step] = addmod(a[i+j+step], a[i+j])
				} else {
					a[i+j+step] = submod(a[i+j+step], a[i+j])
				}
			}
		}
	}
}

func fwtAnd(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				if !invert {
					a[i+j] = addmod(a[i+j], a[i+j+step])
				} else {
					a[i+j] = submod(a[i+j], a[i+j+step])
				}
			}
		}
	}
}

func powMod(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	freq := make([]int64, SIZE)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		freq[v]++
	}

	fib := make([]int64, SIZE)
	fib[1] = 1
	for i := 2; i < SIZE; i++ {
		fib[i] = addmod(fib[i-1], fib[i-2])
	}

	// pair counts with OR
	a := make([]int64, SIZE)
	copy(a, freq)
	fwtOr(a, false)
	for i := 0; i < SIZE; i++ {
		a[i] = a[i] * a[i] % MOD
	}
	fwtOr(a, true)

	// single counts
	b := make([]int64, SIZE)
	copy(b, freq)

	// pair counts with XOR
	c := make([]int64, SIZE)
	copy(c, freq)
	fwtXor(c, false)
	for i := 0; i < SIZE; i++ {
		c[i] = c[i] * c[i] % MOD
	}
	fwtXor(c, true)

	for i := 0; i < SIZE; i++ {
		a[i] = a[i] * fib[i] % MOD
		b[i] = b[i] * fib[i] % MOD
		c[i] = c[i] * fib[i] % MOD
	}

	fwtAnd(a, false)
	fwtAnd(b, false)
	fwtAnd(c, false)

	for i := 0; i < SIZE; i++ {
		a[i] = a[i] * b[i] % MOD * c[i] % MOD
	}

	fwtAnd(a, true)

	var ans int64
	for i := 0; i < BITS; i++ {
		ans = (ans + a[1<<i]) % MOD
	}

	fmt.Println(ans)
}
