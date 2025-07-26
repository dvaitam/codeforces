package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func count(n int, x int64, z int64, choose []int64) int64 {
	if x < 0 {
		return 0
	}
	maxBit := 61
	dp := make([][2]int64, n+1)
	dp[0][0] = 1
	for i := 0; i < maxBit; i++ {
		xi := int((x >> i) & 1)
		zi := int((z >> i) & 1)
		next := make([][2]int64, n+1)
		for c := 0; c <= n; c++ {
			for less := 0; less < 2; less++ {
				ways := dp[c][less]
				if ways == 0 {
					continue
				}
				bitSum := (c & 1) ^ zi
				if less == 0 && bitSum > xi {
					continue
				}
				start := (c - bitSum + 1) / 2
				if start < 0 {
					start = 0
				}
				end := (c + n - bitSum) / 2
				if end > n {
					end = n
				}
				newLessBase := less
				if less == 0 && bitSum < xi {
					newLessBase = 1
				}
				for cp := start; cp <= end; cp++ {
					k := 2*cp + bitSum - c
					if k < 0 || k > n {
						continue
					}
					val := ways * choose[k] % MOD
					next[cp][newLessBase] = (next[cp][newLessBase] + val) % MOD
				}
			}
		}
		dp = next
	}
	return (dp[0][0] + dp[0][1]) % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	var l, r, z int64
	fmt.Fscan(reader, &n, &l, &r, &z)
	choose := make([]int64, n+1)
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	for k := 0; k <= n; k++ {
		choose[k] = fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}
	ans := count(n, r, z, choose) - count(n, l-1, z, choose)
	ans %= MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(writer, ans)
}
