package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const negInf int64 = -1 << 63 / 4
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = negInf
	}
	for _, v := range a {
		next := make([]int64, n+1)
		copy(next, dp)
		for j := 0; j < n; j++ {
			if dp[j] == negInf {
				continue
			}
			val := dp[j] + int64(j+1)*v
			if val > next[j+1] {
				next[j+1] = val
			}
		}
		dp = next
	}
	ans := negInf
	for _, v := range dp {
		if v > ans {
			ans = v
		}
	}
	if ans == negInf {
		ans = 0
	}
	fmt.Println(ans)
}
