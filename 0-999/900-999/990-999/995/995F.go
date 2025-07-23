package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var D int64
	fmt.Fscan(in, &n, &D)
	parent := make([]int, n+1)
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent[i])
		children[parent[i]] = append(children[parent[i]], i)
	}

	m := n
	dp := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]int64, m+1)
	}

	order := make([]int, 0, n)
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, c := range children[v] {
			stack = append(stack, c)
		}
	}

	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		for x := 1; x <= m; x++ {
			prod := int64(1)
			for _, c := range children[v] {
				prod = prod * dp[c][x] % MOD
			}
			dp[v][x] = (dp[v][x-1] + prod) % MOD
		}
	}

	f := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		f[i] = dp[1][i]
	}

	if D <= int64(m) {
		fmt.Println(f[D])
		return
	}

	fact := make([]int64, m+1)
	invFact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[m] = modPow(fact[m], MOD-2)
	for i := m - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}

	pre := make([]int64, m+1)
	suf := make([]int64, m+2)
	pre[0] = 1
	for i := 1; i <= m; i++ {
		val := (D - int64(i-1)) % MOD
		if val < 0 {
			val += MOD
		}
		pre[i] = pre[i-1] * val % MOD
	}
	suf[m+1] = 1
	for i := m; i >= 0; i-- {
		val := (D - int64(i)) % MOD
		if val < 0 {
			val += MOD
		}
		suf[i] = suf[i+1] * val % MOD
	}

	ans := int64(0)
	for i := 0; i <= m; i++ {
		num := pre[i] * suf[i+1] % MOD
		term := f[i] * num % MOD
		term = term * invFact[i] % MOD
		term = term * invFact[m-i] % MOD
		if (m-i)%2 == 1 {
			term = (MOD - term) % MOD
		}
		ans = (ans + term) % MOD
	}

	fmt.Println(ans)
}
