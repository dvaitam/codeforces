package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		parent := make([]int, n+1)
		depth := make([]int, n+1)
		maxDepth := 0
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &parent[i])
			depth[i] = depth[parent[i]] + 1
			if depth[i] > maxDepth {
				maxDepth = depth[i]
			}
		}

		levels := make([][]int, maxDepth+1)
		for i := 1; i <= n; i++ {
			d := depth[i]
			if d >= len(levels) {
				tmp := make([][]int, d+1)
				copy(tmp, levels)
				levels = tmp
			}
			levels[d] = append(levels[d], i)
		}

		dp := make([]int64, n+1)
		S := make([]int64, maxDepth+1)

		dp[1] = 1
		S[0] = 1

		if maxDepth >= 1 {
			for _, v := range levels[1] {
				if v == 1 {
					continue
				}
				dp[v] = 1
				S[1] = (S[1] + 1) % MOD
			}
		}

		for d := 2; d <= maxDepth; d++ {
			sumPrev := S[d-1]
			for _, v := range levels[d] {
				p := parent[v]
				val := (sumPrev - dp[p]) % MOD
				if val < 0 {
					val += MOD
				}
				dp[v] = val
				S[d] = (S[d] + val) % MOD
			}
		}

		var ans int64
		for d := 0; d <= maxDepth; d++ {
			ans = (ans + S[d]) % MOD
		}
		fmt.Fprintln(out, ans)
	}
}
