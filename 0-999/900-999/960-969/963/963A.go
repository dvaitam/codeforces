package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000009

func powmod(a, e int64) int64 {
	a %= mod
	if a < 0 {
		a += mod
	}
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

	var n int64
	var a, b int64
	var k int
	if _, err := fmt.Fscan(reader, &n, &a, &b, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	invA := powmod(a, mod-2)
	cur := powmod(a, n)
	bPow := int64(1)
	sum := int64(0)
	for i := 0; i < k; i++ {
		sign := int64(1)
		if s[i] == '-' {
			sign = -1
		}
		val := cur * bPow % mod
		sum = (sum + sign*val) % mod
		cur = cur * invA % mod
		bPow = bPow * b % mod
	}
	if sum < 0 {
		sum += mod
	}

	ratio := b % mod * invA % mod
	r := powmod(ratio, int64(k))
	m := (n + 1) / int64(k)

	var factor int64
	if r == 1 {
		factor = m % mod
	} else {
		numerator := (powmod(r, m) - 1 + mod) % mod
		denominator := (r - 1 + mod) % mod
		factor = numerator * powmod(denominator, mod-2) % mod
	}

	ans := sum * factor % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
