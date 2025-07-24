package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if n%2 == 1 {
		fmt.Fprintln(writer, 0)
		return
	}

	fac := make([]int64, n+1)
	ifac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[n] = modPow(fac[n], mod-2)
	for i := n; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}

	ans := int64(0)
	for s := 1; s <= n-1; s++ {
		comb := fac[n-2]
		comb = comb * ifac[s-1] % mod
		comb = comb * ifac[n-1-s] % mod
		term := comb
		term = term * modPow(int64(s), int64(s-1)) % mod
		term = term * modPow(int64(n-s), int64(n-s-1)) % mod
		if s%2 == 1 {
			ans = (ans - term) % mod
		} else {
			ans = (ans + term) % mod
		}
	}
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
