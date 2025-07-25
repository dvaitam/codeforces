package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}
	_ = M
	var A, B string
	fmt.Fscan(in, &A)
	fmt.Fscan(in, &B)

	freqA := make([]int, 26)
	freqB := make([]int, 27) // extra for virtual letter after 'Z'
	for _, ch := range A {
		freqA[int(ch-'A')]++
	}
	for _, ch := range B {
		freqB[int(ch-'A')]++
	}

	fact := make([]int64, N+1)
	invFact := make([]int64, N+1)
	fact[0] = 1
	for i := 1; i <= N; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[N] = modPow(fact[N], MOD-2)
	for i := N; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	dp := []int64{1}
	for i := 0; i < 26; i++ {
		newDP := make([]int64, freqA[i]+1)
		prefix := make([]int64, len(dp))
		sum := int64(0)
		for j, v := range dp {
			sum += v
			if sum >= MOD {
				sum -= MOD
			}
			prefix[j] = sum
		}
		for newk := 0; newk <= freqA[i]; newk++ {
			limit := freqB[i] - freqA[i] + newk
			var val int64
			if limit >= 0 {
				if limit >= len(dp) {
					val = prefix[len(dp)-1]
				} else {
					val = prefix[limit]
				}
			}
			val = val * invFact[newk] % MOD
			val = val * invFact[freqA[i]-newk] % MOD
			newDP[newk] = val
		}
		dp = newDP
	}

	ans := fact[N] * dp[0] % MOD
	fmt.Fprintln(out, ans)
}
