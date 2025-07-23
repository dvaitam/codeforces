package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007
const modExp int64 = mod - 1

func powMod(base, exp int64) int64 {
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
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	if k == -1 && (n%2 != m%2) {
		fmt.Println(0)
		return
	}
	exp := ((n - 1) % modExp * ((m - 1) % modExp)) % modExp
	fmt.Println(powMod(2, exp))
}
