package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

var pow10 []int
var C [][]int
var inv9 int

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func precompute(n int) {
	pow10 = make([]int, n+10)
	pow10[0] = 1
	for i := 1; i < len(pow10); i++ {
		pow10[i] = pow10[i-1] * 10 % MOD
	}
	C = make([][]int, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int, i+1)
		C[i][0], C[i][i] = 1, 1
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
		}
	}
	inv9 = modPow(9, MOD-2)
}

type key struct {
	R int
	C [10]int
}

var memo map[key]int

func sumWithPrefix(cnt [10]int, R int) int {
	k := key{R: R, C: cnt}
	if val, ok := memo[k]; ok {
		return val
	}
	prefixGreater := [10]int{}
	for d := 8; d >= 0; d-- {
		prefixGreater[d] = prefixGreater[d+1] + cnt[d+1]
	}
	ways := make([]int, R+1)
	sums := make([]int, R+1)
	ways[0] = 1
	for d := 9; d >= 1; d-- {
		pre := cnt[d]
		greater := prefixGreater[d]
		newWays := make([]int, R+1)
		newSums := make([]int, R+1)
		for t := 0; t <= R; t++ {
			if ways[t] == 0 && sums[t] == 0 {
				continue
			}
			rem := R - t
			for add := 0; add <= rem; add++ {
				nk := t + add
				choose := C[rem][add]
				total := pre + add
				contrib := 0
				if total > 0 {
					contrib = d * pow10[greater+t] % MOD
					contrib = contrib * (pow10[total] - 1 + MOD) % MOD
					contrib = contrib * inv9 % MOD
				}
				newSums[nk] = (newSums[nk] + sums[t]*choose + ways[t]*choose%MOD*contrib) % MOD
				newWays[nk] = (newWays[nk] + ways[t]*choose) % MOD
			}
		}
		ways, sums = newWays, newSums
	}
	total := 0
	for _, v := range sums {
		total = (total + v) % MOD
	}
	memo[k] = total
	return total
}

func F(L int) int {
	var cnt [10]int
	return sumWithPrefix(cnt, L)
}

func solve(X string) int {
	L := len(X)
	preCounts := [10]int{}
	ans := 0
	if L > 0 {
		ans = F(L - 1)
	}
	for i := 0; i < L; i++ {
		digit := int(X[i] - '0')
		for d := 0; d < digit; d++ {
			if i == 0 && d == 0 {
				continue
			}
			preCounts[d]++
			ans = (ans + sumWithPrefix(preCounts, L-i-1)) % MOD
			preCounts[d]--
		}
		preCounts[digit]++
	}
	ans = (ans + sumWithPrefix(preCounts, 0)) % MOD
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var X string
	if _, err := fmt.Fscan(in, &X); err != nil {
		return
	}
	precompute(len(X))
	memo = make(map[key]int)
	fmt.Println(solve(X))
}
