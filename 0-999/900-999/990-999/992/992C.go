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
	var x, k int64
	if _, err := fmt.Fscan(in, &x, &k); err != nil {
		return
	}
	if x == 0 {
		fmt.Println(0)
		return
	}
	pow2k := modPow(2, k)
	pow2k1 := pow2k * 2 % mod
	ans := (pow2k1*(x%mod)%mod - pow2k + 1) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
