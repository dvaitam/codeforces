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

	dp := make([]int64, n+1)
	dp[0] = 1
	ans := make([]int64, n+1)

	for step := 0; ; step++ {
		base := k + step
		if base > n {
			break
		}
		newdp := make([]int64, n+1)
		for r := 0; r < base; r++ {
			sum := int64(0)
			for x := r; x+base <= n; x += base {
				sum += dp[x]
				if sum >= mod {
					sum -= mod
				}
				newdp[x+base] = sum
			}
		}
		hasNonZero := false
		for i := base; i <= n; i++ {
			if newdp[i] != 0 {
				hasNonZero = true
			}
			ans[i] += newdp[i]
			if ans[i] >= mod {
				ans[i] -= mod
			}
		}
		if !hasNonZero {
			break
		}
		dp = newdp
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
}
