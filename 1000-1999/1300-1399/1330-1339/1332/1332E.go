package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	a %= mod
	res := int64(1)
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
	reader := bufio.NewReader(os.Stdin)
	var n, m, L, R int64
	if _, err := fmt.Fscan(reader, &n, &m, &L, &R); err != nil {
		return
	}
	N := n * m
	total := R - L + 1

	if N%2 == 1 {
		fmt.Println(modPow(total%mod, N))
		return
	}

	even := R/2 - (L-1)/2
	odd := total - even
	diff := (even - odd) % mod
	if diff < 0 {
		diff += mod
	}

	a := modPow(total%mod, N)
	b := modPow(diff, N)
	ans := (a + b) % mod
	ans = ans * ((mod + 1) / 2) % mod
	fmt.Println(ans)
}
