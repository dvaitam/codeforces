package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var W, B, X int64
	if _, err := fmt.Fscan(in, &n, &W, &B, &X); err != nil {
		return
	}
	c := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
		sum += c[i]
	}
	cost := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cost[i])
	}

	const neg = int64(-1 << 60)
	dp := make([]int64, sum+1)
	for i := 1; i <= sum; i++ {
		dp[i] = neg
	}
	dp[0] = W

	maxBirds := 0
	for i := 0; i < n; i++ {
		if i > 0 {
			for j := 0; j <= maxBirds; j++ {
				if dp[j] < 0 {
					continue
				}
				cap := W + int64(j)*B
				val := dp[j] + X
				if val > cap {
					val = cap
				}
				dp[j] = val
			}
		}

		newdp := make([]int64, sum+1)
		for k := 0; k <= sum; k++ {
			newdp[k] = neg
		}
		for j := 0; j <= maxBirds; j++ {
			if dp[j] < 0 {
				continue
			}
			maxT := c[i]
			if maxT+j > sum {
				maxT = sum - j
			}
			for t := 0; t <= maxT; t++ {
				req := int64(t) * cost[i]
				if req > dp[j] {
					break
				}
				val := dp[j] - req
				if val > newdp[j+t] {
					newdp[j+t] = val
				}
			}
		}
		maxBirds += c[i]
		dp = newdp
	}

	ans := 0
	for j := maxBirds; j >= 0; j-- {
		if dp[j] >= 0 {
			ans = j
			break
		}
	}
	fmt.Fprintln(out, ans)
}
