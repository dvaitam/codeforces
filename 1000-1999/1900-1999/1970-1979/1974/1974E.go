package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func solve(m int, x int64, c, h []int) int {
	maxH := 0
	for _, v := range h {
		maxH += v
	}
	const inf int64 = math.MaxInt64 / 4
	dp := make([]int64, maxH+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for i := 1; i <= m; i++ {
		for j := maxH; j >= 0; j-- {
			if dp[j] == inf {
				continue
			}
			if int64(j) > int64(maxH) {
				continue
			}
			if dp[j]+int64(c[i-1]) <= int64(i-1)*x {
				if nj := j + h[i-1]; nj <= maxH {
					if val := dp[j] + int64(c[i-1]); val < dp[nj] {
						dp[nj] = val
					}
				}
			}
		}
	}
	res := 0
	for i := maxH; i >= 0; i-- {
		if dp[i] != inf {
			res = i
			break
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var m int
		var x int64
		fmt.Fscan(in, &m, &x)
		c := make([]int, m)
		h := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &c[i], &h[i])
		}
		fmt.Fprintln(out, solve(m, x, c, h))
	}
}
