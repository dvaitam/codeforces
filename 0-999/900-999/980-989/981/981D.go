package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(prefix []int64, n, k int, mask int64) bool {
	dp := make([]bool, n+1)
	dp[0] = true
	for step := 0; step < k; step++ {
		next := make([]bool, n+1)
		for i := 0; i <= n; i++ {
			if !dp[i] {
				continue
			}
			for j := i + 1; j <= n; j++ {
				if ((prefix[j] - prefix[i]) & mask) == mask {
					next[j] = true
				}
			}
		}
		dp = next
	}
	return dp[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}
	var ans int64
	for bit := 60; bit >= 0; bit-- {
		candidate := ans | (1 << uint(bit))
		if can(prefix, n, k, candidate) {
			ans = candidate
		}
	}
	fmt.Println(ans)
}
