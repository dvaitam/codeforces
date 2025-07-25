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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	const INF int64 = 1 << 60
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}
		dp := [3]int64{0, 0, 0}
		for k := 0; k < 3; k++ {
			dp[k] = int64(k) * b[0]
		}
		for i := 1; i < n; i++ {
			ndp := [3]int64{INF, INF, INF}
			for k := 0; k < 3; k++ {
				for j := 0; j < 3; j++ {
					if a[i]+int64(k) != a[i-1]+int64(j) {
						val := dp[j] + int64(k)*b[i]
						if val < ndp[k] {
							ndp[k] = val
						}
					}
				}
			}
			dp = ndp
		}
		ans := dp[0]
		if dp[1] < ans {
			ans = dp[1]
		}
		if dp[2] < ans {
			ans = dp[2]
		}
		fmt.Fprintln(out, ans)
	}
}
