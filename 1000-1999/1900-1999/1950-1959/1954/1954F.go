package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = int(int64(res) * int64(a) % int64(MOD))
		}
		a = int(int64(a) * int64(a) % int64(MOD))
		b >>= 1
	}
	return res
}

func phi(x int) int {
	res := x
	i := 2
	for i*i <= x {
		if x%i == 0 {
			for x%i == 0 {
				x /= i
			}
			res -= res / i
		}
		i++
	}
	if x > 1 {
		res -= res / x
	}
	return res
}

func F(n, c, k, d int) int {
	cost := make([]int, d)
	for i := c; i < n; i++ {
		cost[i%d]++
	}
	mandatory := make([]bool, d)
	mandatoryCost := 0
	for i := 0; i < c; i++ {
		r := i % d
		if !mandatory[r] {
			mandatory[r] = true
			mandatoryCost += cost[r]
		}
	}
	if mandatoryCost > k {
		return 0
	}
	limit := k - mandatoryCost
	dp := make([]int, limit+1)
	dp[0] = 1
	for r := 0; r < d; r++ {
		if mandatory[r] {
			continue
		}
		w := cost[r]
		for t := limit; t >= w; t-- {
			dp[t] += dp[t-w]
			if dp[t] >= MOD {
				dp[t] -= MOD
			}
		}
	}
	total := 0
	for t := 0; t <= limit; t++ {
		total += dp[t]
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

	var n, c, k int
	if _, err := fmt.Fscan(reader, &n, &c, &k); err != nil {
		return
	}

	divisors := []int{}
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			divisors = append(divisors, d)
			if d*d != n {
				divisors = append(divisors, n/d)
			}
		}
	}

	ans := 0
	for _, d := range divisors {
		f := F(n, c, k, d)
		phiVal := phi(n / d)
		ans = (ans + int(int64(f)*int64(phiVal)%int64(MOD))) % MOD
	}

	invN := modPow(n%MOD, MOD-2)
	ans = int(int64(ans) * int64(invN) % int64(MOD))
	fmt.Fprintln(writer, ans)
}
