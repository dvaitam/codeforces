package main

import (
	"fmt"
)

const mod int64 = 1000000007

func powmod(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	total := powmod(3, 3*n)
	bad := powmod(7, n)
	ans := total - bad
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
