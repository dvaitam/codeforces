package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, T int
	if _, err := fmt.Fscan(reader, &n, &T); err != nil {
		return
	}
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}
	// dp[2][T+1]
	dp := make([][]float64, 2)
	dp[0] = make([]float64, T+1)
	dp[1] = make([]float64, T+1)
	cur := 0
	// iterate in reverse order
	for idx := n - 1; idx >= 0; idx-- {
		cur ^= 1
		prev := cur ^ 1
		p := float64(x[idx]) / 100.0
		// P = (1-p)^(y[idx]-1)
		P := math.Pow(1-p, float64(y[idx]-1))
		sum := 0.0
		dp[cur][0] = 0.0
		for j := 1; j <= T; j++ {
			sum *= 1 - p
			sum += (dp[prev][j-1] + 1.0) * p
			if j >= y[idx] {
				sum += (dp[prev][j-y[idx]] + 1.0) * P * (1 - p)
			}
			if j > y[idx] {
				sum -= (dp[prev][j-y[idx]-1] + 1.0) * P * (1 - p)
			}
			dp[cur][j] = sum
		}
	}
	result := dp[cur][T]
	// print with 10 decimal places
	fmt.Printf("%.10f\n", result)
}
