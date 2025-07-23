package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000003

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

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, Cval int
	if _, err := fmt.Fscan(reader, &n, &Cval); err != nil {
		return
	}

	maxN := n + Cval
	fact, inv := prepareFact(maxN)

	res := (C(int64(maxN), int64(Cval), fact, inv) - 1 + mod) % mod
	fmt.Fprintln(writer, res)
}
