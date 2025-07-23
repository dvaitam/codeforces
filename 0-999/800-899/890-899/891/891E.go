package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		a[i] %= mod
	}

	e := make([]int64, n+1)
	e[0] = 1
	for _, v := range a {
		for j := n; j >= 1; j-- {
			e[j] = (e[j] + e[j-1]*v) % mod
		}
	}

	invN := modPow(int64(n), mod-2)
	limit := n
	if int64(limit) > k {
		limit = int(k)
	}
	F := int64(1)
	Ep := int64(0)
	for t := 0; t <= limit; t++ {
		term := e[n-t] * F % mod
		if t%2 == 1 {
			term = (mod - term) % mod
		}
		Ep = (Ep + term) % mod
		if t < limit {
			mult := (k - int64(t)) % mod
			F = F * mult % mod * invN % mod
		}
	}
	ans := (e[n] - Ep) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
