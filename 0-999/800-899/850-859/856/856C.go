package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

var fact, inv []int64

func powMod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initFact(n int) {
	fact = make([]int64, n+1)
	inv = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv[n] = powMod(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * inv[k] % MOD * inv[n-k] % MOD
}

func subsetDP(vals []int) [][]int64 {
	m := len(vals)
	dp := make([][]int64, m+1)
	for i := range dp {
		dp[i] = make([]int64, 11)
	}
	dp[0][0] = 1
	for _, v := range vals {
		for i := m - 1; i >= 0; i-- {
			for r := 0; r < 11; r++ {
				if dp[i][r] != 0 {
					dp[i+1][(r+v)%11] = (dp[i+1][(r+v)%11] + dp[i][r]) % MOD
				}
			}
		}
	}
	return dp
}

func digitLen(x int) int {
	l := 0
	for x > 0 {
		l++
		x /= 10
	}
	if l == 0 {
		l = 1
	}
	return l
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)

	maxN := 4000
	initFact(maxN)

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &nums[i])
		}

		oddVals := make([]int, 0)
		evenVals := make([]int, 0)
		for _, x := range nums {
			l := digitLen(x)
			m := x % 11
			if l%2 == 1 {
				oddVals = append(oddVals, m)
			} else {
				evenVals = append(evenVals, m)
			}
		}

		n1 := len(oddVals)
		n0 := len(evenVals)
		sumOdd := 0
		for _, v := range oddVals {
			sumOdd = (sumOdd + v) % 11
		}
		sumEven := 0
		for _, v := range evenVals {
			sumEven = (sumEven + v) % 11
		}

		dpOdd := subsetDP(oddVals)
		kneg := (n1 + 1) / 2
		pos := n1 - kneg
		oddWays := make([]int64, 11)
		for r := 0; r < 11; r++ {
			cnt := dpOdd[kneg][r]
			if cnt != 0 {
				t := (sumOdd - 2*r) % 11
				if t < 0 {
					t += 11
				}
				oddWays[t] = (oddWays[t] + cnt*fact[kneg]%MOD*fact[pos]) % MOD
			}
		}

		dpEven := subsetDP(evenVals)
		P := n1/2 + 1
		M := n1 + 1 - P
		evenWays := make([]int64, 11)
		for k := 0; k <= n0; k++ {
			if M == 0 {
				if n0-k > 0 {
					continue
				}
				coef := comb(k+P-1, P-1)
				coef = coef * fact[k] % MOD * fact[n0-k] % MOD
				for r := 0; r < 11; r++ {
					cnt := dpEven[k][r]
					if cnt == 0 {
						continue
					}
					u := (2*r - sumEven) % 11
					if u < 0 {
						u += 11
					}
					evenWays[u] = (evenWays[u] + cnt*coef) % MOD
				}
			} else {
				coef := comb(k+P-1, P-1) * comb(n0-k+M-1, M-1) % MOD
				coef = coef * fact[k] % MOD * fact[n0-k] % MOD
				for r := 0; r < 11; r++ {
					cnt := dpEven[k][r]
					if cnt == 0 {
						continue
					}
					u := (2*r - sumEven) % 11
					if u < 0 {
						u += 11
					}
					evenWays[u] = (evenWays[u] + cnt*coef) % MOD
				}
			}
		}

		var ans int64
		for t := 0; t < 11; t++ {
			ans = (ans + oddWays[t]*evenWays[(11-t)%11]) % MOD
		}
		fmt.Fprintln(writer, ans)
	}
}
