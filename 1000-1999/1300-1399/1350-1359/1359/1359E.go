package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func powMod(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initFactorials(n int) ([]int, []int) {
	fact := make([]int, n+1)
	inv := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	inv[n] = powMod(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * i % MOD
	}
	return fact, inv
}

func comb(n, k int, fact, inv []int) int {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * inv[k] % MOD * inv[n-k] % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	fact, inv := initFactorials(n)
	ans := 0
	for m := 1; m <= n; m++ {
		cnt := n / m
		if cnt >= k {
			ans += comb(cnt-1, k-1, fact, inv)
			if ans >= MOD {
				ans -= MOD
			}
		}
	}
	fmt.Println(ans)
}
