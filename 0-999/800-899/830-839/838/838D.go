package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(base, exp int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	if m > n {
		fmt.Println(0)
		return
	}
	pow1 := modPow(n+1, m-1)
	pow2 := modPow(2, m)
	ans := (n - m + 1) % mod
	ans = ans * pow1 % mod
	ans = ans * pow2 % mod
	fmt.Println(ans)
}
