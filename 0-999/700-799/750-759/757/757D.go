package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

type Edge struct {
	to  int
	val int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	edges := make([][]Edge, n)
	for i := 0; i < n; i++ {
		val := 0
		for j := i; j < n; j++ {
			val = val*2 + int(s[j]-'0')
			if val > 20 {
				break
			}
			if val > 0 {
				edges[i] = append(edges[i], Edge{j + 1, val})
			}
		}
	}

	dp := make([]map[int]int64, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make(map[int]int64)
	}
	for i := 0; i < n; i++ {
		dp[i][0] = (dp[i][0] + 1) % MOD
	}

	for i := 0; i < n; i++ {
		for mask, cnt := range dp[i] {
			for _, e := range edges[i] {
				nm := mask | (1 << (e.val - 1))
				dp[e.to][nm] = (dp[e.to][nm] + cnt) % MOD
			}
		}
	}

	var ans int64
	for i := 0; i <= n; i++ {
		for mask, cnt := range dp[i] {
			if mask != 0 && (mask&(mask+1)) == 0 {
				ans = (ans + cnt) % MOD
			}
		}
	}

	fmt.Fprintln(out, ans)
}
