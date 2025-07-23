package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	two := make([]int, n)
	five := make([]int, n)
	maxFive := 0
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		c2, c5 := 0, 0
		for x%2 == 0 {
			c2++
			x /= 2
		}
		for x%5 == 0 {
			c5++
			x /= 5
		}
		two[i] = c2
		five[i] = c5
		maxFive += c5
	}
	dp := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int, maxFive+1)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	dp[0][0] = 0
	for idx := 0; idx < n; idx++ {
		t2, f5 := two[idx], five[idx]
		for j := k; j >= 1; j-- {
			for f := maxFive; f >= f5; f-- {
				if dp[j-1][f-f5] >= 0 {
					val := dp[j-1][f-f5] + t2
					if val > dp[j][f] {
						dp[j][f] = val
					}
				}
			}
		}
	}
	ans := 0
	for f := 0; f <= maxFive; f++ {
		val := dp[k][f]
		if val < 0 {
			continue
		}
		if val < f {
			if val > ans {
				ans = val
			}
		} else {
			if f > ans {
				ans = f
			}
		}
	}
	fmt.Println(ans)
}
