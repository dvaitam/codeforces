package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

var fact, invFact []int64

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func initComb(n int) {
	fact = make([]int64, n+1)
	invFact = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)

	maxN := 100000
	initComb(maxN)

	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)

		zeros := 0
		pairs := 0
		for i := 0; i < n; {
			if s[i] == '0' {
				zeros++
				i++
			} else {
				j := i
				for j < n && s[j] == '1' {
					j++
				}
				length := j - i
				pairs += length / 2
				i = j
			}
		}
		ans := comb(zeros+pairs, pairs)
		fmt.Fprintln(writer, ans)
	}
}
