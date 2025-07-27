package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k, r, c int64
	if _, err := fmt.Fscan(in, &n, &m, &k, &r, &c); err != nil {
		return
	}
	var ax, ay, bx, by int64
	fmt.Fscan(in, &ax, &ay, &bx, &by)

	var exp int64
	if ax == bx && ay == by {
		exp = n * m
	} else {
		exp = n*m - r*c
	}

	fmt.Println(modPow(k%mod, exp))
}
