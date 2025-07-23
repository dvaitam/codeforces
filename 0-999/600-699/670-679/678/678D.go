package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powmod(a, e int64) int64 {
	a %= mod
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

func inv(a int64) int64 {
	return powmod(a, mod-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var A, B, n, x int64
	if _, err := fmt.Fscan(in, &A, &B, &n, &x); err != nil {
		return
	}
	if A == 1 {
		ans := (x%mod + (B%mod)*(n%mod)) % mod
		fmt.Println(ans)
		return
	}
	powA := powmod(A, n)
	term := (powA - 1 + mod) % mod
	term = term * (B % mod) % mod
	term = term * inv((A-1)%mod) % mod
	ans := (powA*(x%mod) + term) % mod
	fmt.Println(ans)
}
