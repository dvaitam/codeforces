package main

import (
	"bufio"
	"fmt"
	"os"
)

var bases = []uint64{2, 3, 5, 7, 11}

func modMul(a, b, mod uint64) uint64 {
	return (a * b) % mod
}

func modPow(a, d, mod uint64) uint64 {
	res := uint64(1)
	for d > 0 {
		if d&1 == 1 {
			res = modMul(res, a%mod, mod)
		}
		a = modMul(a%mod, a%mod, mod)
		d >>= 1
	}
	return res
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	smallPrimes := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for _, p := range smallPrimes {
		if n == p {
			return true
		}
		if n%p == 0 {
			return false
		}
	}
	d := n - 1
	r := 0
	for d&1 == 0 {
		d >>= 1
		r++
	}
	for _, a := range bases {
		if a%(n) == 0 {
			continue
		}
		x := modPow(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for i := 1; i < r; i++ {
			x = modMul(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x uint64
		var k int
		fmt.Fscan(in, &x, &k)
		if k == 1 && isPrime(x) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
