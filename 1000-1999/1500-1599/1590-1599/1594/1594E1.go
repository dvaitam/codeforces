package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

func powMod(base, exp int64) int64 {
	base %= mod
	res := int64(1)
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
	reader := bufio.NewReader(os.Stdin)
	var k int64
	fmt.Fscan(reader, &k)

	exp := (int64(1) << (k + 1)) - 3
	ans := powMod(2, exp) * 3 % mod
	fmt.Println(ans)
}
