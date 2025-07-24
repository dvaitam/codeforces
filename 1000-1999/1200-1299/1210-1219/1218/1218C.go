package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var N, M, K int
	if _, err := fmt.Fscan(in, &N, &M, &K); err != nil {
		return
	}
	cost := make([][]int64, N)
	for i := range cost {
		cost[i] = make([]int64, M)
	}
	for i := 0; i < K; i++ {
		var x, y, d, t int
		var e int64
		fmt.Fscan(in, &x, &y, &d, &t, &e)
		// first pos (x,y) at times t + 4k
		timeA := x + y
		if timeA >= t && (timeA-t)%4 == 0 {
			cost[x][y] += e
		}
		// second pos1 (x+d, y-d) at times t+1 + 4k
		if timeA >= t+1 && (timeA-(t+1))%4 == 0 {
			cost[x+d][y-d] += e
		}
		timeB := x + y + d
		// third pos (x+d,y) at t+2+4k
		if timeB >= t+2 && (timeB-(t+2))%4 == 0 {
			cost[x+d][y] += e
		}
		// fourth pos (x,y+d) at t+3+4k
		if timeB >= t+3 && (timeB-(t+3))%4 == 0 {
			cost[x][y+d] += e
		}
	}
	dp := make([][]int64, N)
	for i := range dp {
		dp[i] = make([]int64, M)
	}
	dp[0][0] = cost[0][0]
	for i := 1; i < M; i++ {
		dp[0][i] = dp[0][i-1] + cost[0][i]
	}
	for i := 1; i < N; i++ {
		dp[i][0] = dp[i-1][0] + cost[i][0]
	}
	for i := 1; i < N; i++ {
		for j := 1; j < M; j++ {
			if dp[i-1][j] < dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
			dp[i][j] += cost[i][j]
		}
	}
	fmt.Println(dp[N-1][M-1])
}
