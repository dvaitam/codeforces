package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func modPow(a, b int64) int64 {
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
	in := bufio.NewReader(os.Stdin)
	var n int64
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := int64(0); i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if k == 1 {
		fmt.Println(0)
		return
	}
	// Placeholder formula: (k^n - k)/(k-1)
	kn := modPow(k%mod, n)
	ans := (kn - k%mod + mod) % mod
	inv := modPow(k-1, mod-2)
	ans = ans * inv % mod
	fmt.Println(ans)
}
