package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

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

	positions := make([][]int, 26)
	for i := 0; i < n; i++ {
		c := int(s[i] - 'a')
		positions[c] = append(positions[c], i)
	}

	dp := make([]int64, n+1)
	dp[0] = 1
	for step := 0; step < n; step++ {
		prefix := make([]int64, n+1)
		cur := int64(0)
		for i := 0; i < n; i++ {
			cur += dp[i]
			if cur >= mod {
				cur -= mod
			}
			prefix[i] = cur
		}
		newDp := make([]int64, n+1)
		for _, posList := range positions {
			if len(posList) == 0 {
				continue
			}
			prev := -1
			for _, p := range posList {
				sum := prefix[p]
				if prev >= 0 {
					sum -= prefix[prev]
					if sum < 0 {
						sum += mod
					}
				}
				newDp[p] += sum
				if newDp[p] >= mod {
					newDp[p] %= mod
				}
				prev = p
			}
		}
		dp = newDp
	}

	ans := int64(0)
	for i := 0; i < n; i++ {
		ans += dp[i]
		if ans >= mod {
			ans -= mod
		}
	}
	fmt.Fprintln(out, ans%mod)
}
