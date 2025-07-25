package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &h[i])
	}

	offset := n
	dp := make([]int64, 2*n+3)
	dp[offset] = 1

	for i := 0; i < n; i++ {
		hi := h[i]
		hi1 := h[(i+1)%n]
		if hi == hi1 {
			mult := int64(k) % mod
			for j := range dp {
				dp[j] = dp[j] * mult % mod
			}
			continue
		}
		next := make([]int64, 2*n+3)
		for j, val := range dp {
			if val == 0 {
				continue
			}
			// choose value equal to hi -> score difference decreases by 1
			if j-1 >= 0 {
				next[j-1] = (next[j-1] + val) % mod
			}
			// choose value equal to hi1 -> score difference increases by 1
			if j+1 < len(next) {
				next[j+1] = (next[j+1] + val) % mod
			}
			// choose any other value -> difference stays the same
			next[j] = (next[j] + val*int64(k-2)) % mod
		}
		dp = next
	}

	var ans int64
	for d := offset + 1; d < len(dp); d++ {
		ans = (ans + dp[d]) % mod
	}
	fmt.Fprintln(writer, ans)
}
