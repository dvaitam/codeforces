package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(a) % int64(mod))
		}
		a = int(int64(a) * int64(a) % int64(mod))
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, c int
	if _, err := fmt.Fscan(reader, &n, &k, &c); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	m := 1 << c

	// precompute factorials and inverse factorials up to k
	fact := make([]int, k+1)
	invFact := make([]int, k+1)
	fact[0] = 1
	for i := 1; i <= k; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % int64(mod))
	}
	invFact[k] = modPow(fact[k], mod-2)
	for i := k; i >= 1; i-- {
		invFact[i-1] = int(int64(invFact[i]) * int64(i) % int64(mod))
	}

	size := (k + 1) * m
	dp := make([]int, size)
	ndp := make([]int, size)
	dp[0] = 1

	for _, val := range a {
		for i := 0; i < size; i++ {
			ndp[i] = 0
		}
		for t := 0; t <= k; t++ {
			base := t * m
			for x := 0; x < m; x++ {
				cur := dp[base+x]
				if cur == 0 {
					continue
				}
				for d := 0; d <= k-t; d++ {
					newT := t + d
					newX := x ^ (val - d)
					idx := newT*m + newX
					ndp[idx] = (ndp[idx] + int(int64(cur)*int64(invFact[d])%int64(mod))) % mod
				}
			}
		}
		dp, ndp = ndp, dp
	}

	nkInv := modPow(modPow(n%mod, k), mod-2)

	for x := 0; x < m; x++ {
		res := dp[k*m+x]
		res = int(int64(res) * int64(fact[k]) % int64(mod))
		res = int(int64(res) * int64(nkInv) % int64(mod))
		if x > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res)
	}
	fmt.Fprintln(writer)
}
