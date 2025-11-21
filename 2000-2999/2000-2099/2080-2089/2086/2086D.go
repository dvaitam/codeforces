package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

type testCase struct {
	counts [26]int
	sum    int
}

func modPow(base, exp int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}

func precompute(maxN int) ([]int64, []int64) {
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	return fact, invFact
}

func solve(tc testCase, fact, invFact []int64) int64 {
	n := tc.sum
	if n == 0 {
		return 0
	}

	odd := (n + 1) / 2
	even := n - odd

	invProd := int64(1)
	counts := make([]int, 0, 26)
	for _, c := range tc.counts {
		invProd = invProd * invFact[c] % mod
		if c > 0 {
			counts = append(counts, c)
		}
	}

	dp := make([]int64, odd+1)
	dp[0] = 1
	for _, c := range counts {
		if c > odd {
			continue
		}
		for s := odd; s >= c; s-- {
			dp[s] += dp[s-c]
			if dp[s] >= mod {
				dp[s] -= mod
			}
		}
	}

	subsetCount := dp[odd]
	if subsetCount == 0 {
		return 0
	}

	ans := subsetCount
	ans = ans * fact[odd] % mod
	ans = ans * fact[even] % mod
	ans = ans * invProd % mod
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	tests := make([]testCase, t)
	maxSum := 0
	for i := 0; i < t; i++ {
		var sum int
		var arr [26]int
		for j := 0; j < 26; j++ {
			fmt.Fscan(in, &arr[j])
			sum += arr[j]
		}
		tests[i] = testCase{counts: arr, sum: sum}
		if sum > maxSum {
			maxSum = sum
		}
	}

	fact, invFact := precompute(maxSum)

	for _, tc := range tests {
		fmt.Fprintln(out, solve(tc, fact, invFact))
	}
}
