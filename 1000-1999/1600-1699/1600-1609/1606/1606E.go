package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

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
	var n, x int64
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	total := modPow(x%mod, n)
	if n == 2 {
		fmt.Println(x % mod)
		return
	}
	if x < n-1 {
		fmt.Println(total)
		return
	}
	winners := n % mod
	winners = winners * modPow(n-1, n-1) % mod
	sub := modPow(2, x-(n-1))
	sub = (sub - 1 + mod) % mod
	winners = winners * sub % mod
	ans := (total - winners) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
