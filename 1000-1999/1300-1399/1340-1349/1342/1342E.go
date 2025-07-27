package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

var fact, invFact []int64

func powmod(a, b int64) int64 {
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

func initFact(n int) {
	fact = make([]int64, n+1)
	invFact = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = powmod(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func C(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * invFact[r] % mod * invFact[n-r] % mod
}

func surjection(n, c int) int64 {
	res := int64(0)
	for i := 0; i <= c; i++ {
		term := C(c, i) * powmod(int64(c-i), int64(n)) % mod
		if i%2 == 1 {
			res = (res - term) % mod
		} else {
			res = (res + term) % mod
		}
	}
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k > n-1 {
		fmt.Fprintln(writer, 0)
		return
	}
	initFact(n)
	if k == 0 {
		fmt.Fprintln(writer, fact[n]%mod)
		return
	}
	c := n - k
	val := C(n, c) * surjection(n, c) % mod
	ans := val * 2 % mod
	fmt.Fprintln(writer, ans)
}
