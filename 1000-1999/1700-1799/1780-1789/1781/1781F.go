package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(a, e int64) int64 {
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

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int64
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	inv10000 := modPow(10000%mod, mod-2)
	p := q % mod * inv10000 % mod
	ans := modPow(p, n)
	fmt.Println(ans)
}
