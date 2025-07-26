package main

import (
	"bufio"
	"fmt"
	"os"
)

var rdr = bufio.NewReader(os.Stdin)
var wrtr = bufio.NewWriter(os.Stdout)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve() {
	var n int
	if _, err := fmt.Fscan(rdr, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := range arr {
		fmt.Fscan(rdr, &arr[i])
	}
	var x int
	fmt.Fscan(rdr, &x)

	const negInf = -1 << 60
	dp := [2][2]int{}
	for i := range dp {
		for j := range dp[i] {
			dp[i][j] = negInf
		}
	}
	dp[0][0] = 0
	for i := 0; i < n; i++ {
		ndp := [2][2]int{}
		for a := range ndp {
			for b := range ndp[a] {
				ndp[a][b] = negInf
			}
		}
		for p2 := 0; p2 < 2; p2++ {
			for p1 := 0; p1 < 2; p1++ {
				cur := dp[p2][p1]
				if cur == negInf {
					continue
				}
				if cur > ndp[p1][0] {
					ndp[p1][0] = cur
				}
				ok := true
				if p1 == 1 && arr[i-1]+arr[i] < 2*x {
					ok = false
				}
				if ok && p1 == 1 && p2 == 1 && i >= 2 && arr[i-2]+arr[i-1]+arr[i] < 3*x {
					ok = false
				}
				if ok {
					if cur+1 > ndp[p1][1] {
						ndp[p1][1] = cur + 1
					}
				}
			}
		}
		dp = ndp
	}
	ans := negInf
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if dp[i][j] > ans {
				ans = dp[i][j]
			}
		}
	}
	fmt.Fprintln(wrtr, ans)
}

func main() {
	defer wrtr.Flush()
	var t int
	fmt.Fscan(rdr, &t)
	for ; t > 0; t-- {
		solve()
	}
}
