package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func mobius(n int) []int {
	mu := make([]int, n+1)
	mu[1] = 1
	primes := []int{}
	isComp := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}
	return mu
}

func countForD(d, n, m int, L, R []int) int64 {
	limit := m / d
	base := 0
	diffs := make([]int, n)
	for i := 0; i < n; i++ {
		li := (L[i] + d - 1) / d
		ri := R[i] / d
		if li > ri {
			return 0
		}
		base += li
		diffs[i] = ri - li
	}
	limit -= base
	if limit < 0 {
		return 0
	}
	dp := make([]int64, limit+1)
	dp[0] = 1
	for _, r := range diffs {
		prefix := int64(0)
		ndp := make([]int64, limit+1)
		for s := 0; s <= limit; s++ {
			prefix += dp[s]
			if prefix >= MOD {
				prefix -= MOD
			}
			if s-r-1 >= 0 {
				prefix -= dp[s-r-1]
				if prefix < 0 {
					prefix += MOD
				}
			}
			ndp[s] = prefix
		}
		dp = ndp
	}
	var total int64
	for _, v := range dp {
		total += v
		if total >= MOD {
			total -= MOD
		}
	}
	return total
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &L[i], &R[i])
	}

	mu := mobius(m)
	var ans int64
	for d := 1; d <= m; d++ {
		if mu[d] == 0 {
			continue
		}
		val := countForD(d, n, m, L, R)
		if val == 0 {
			continue
		}
		if mu[d] == 1 {
			ans += val
		} else {
			ans -= val
		}
		ans %= MOD
	}
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(writer, ans)
}
