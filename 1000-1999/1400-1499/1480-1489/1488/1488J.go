package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func countBouquets(weights []int, counts []int, k int) int {
	dp := make([]int, k+1)
	dp[0] = 1
	for idx, w := range weights {
		c := counts[idx]
		if c == 0 {
			continue
		}
		ndp := make([]int, k+1)
		// process using sliding window for bounded knapsack
		for r := 0; r < w && r <= k; r++ {
			sum := 0
			for j := r; j <= k; j += w {
				sum += dp[j]
				if sum >= mod {
					sum -= mod
				}
				ndp[j] = sum
				if j-w*c >= 0 {
					sum -= dp[j-w*c]
					if sum < 0 {
						sum += mod
					}
				}
			}
		}
		dp = ndp
	}
	ans := 0
	for j := 0; j <= k; j++ {
		ans += dp[j]
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	w := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &w[i])
	}
	stock := make([]int, n)

	for q := 0; q < m; q++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var i, c int
			fmt.Fscan(in, &i, &c)
			stock[i-1] += c
		} else if t == 2 {
			var i, c int
			fmt.Fscan(in, &i, &c)
			if stock[i-1] >= c {
				stock[i-1] -= c
			}
		} else if t == 3 {
			var l, r, k int
			fmt.Fscan(in, &l, &r, &k)
			weights := make([]int, r-l+1)
			counts := make([]int, r-l+1)
			for i := l - 1; i < r; i++ {
				weights[i-(l-1)] = w[i]
				counts[i-(l-1)] = stock[i]
			}
			res := countBouquets(weights, counts, k)
			fmt.Fprintln(out, res)
		}
	}
}
