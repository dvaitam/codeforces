package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n int
	p [][]float64
)

// calc computes the probability based on sequence seq
func calc(seq []int) float64 {
	l := len(seq)
	// reverse seq
	rev := make([]int, l)
	for i, v := range seq {
		rev[l-1-i] = v
	}
	// dp[0] and dp[1] rows
	dp := [2][]float64{}
	dp[0] = make([]float64, n)
	dp[1] = make([]float64, n)
	dp[0][rev[0]] = 1.0
	// iterate transitions
	for i := 0; i < l-1; i++ {
		i0 := i & 1
		i1 := i0 ^ 1
		nw := rev[i+1]
		sum := 0.0
		for j := 0; j < n; j++ {
			dp[i1][j] = dp[i0][j] * p[j][nw]
			sum += dp[i1][j]
		}
		dp[i1][nw] = 1.0 - sum
	}
	return dp[(l-1)&1][0]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	p = make([][]float64, n)
	for i := 0; i < n; i++ {
		p[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &p[i][j])
		}
	}
	use := make([]bool, n)
	seq := make([]int, 1, n)
	seq[0] = 0
	// greedy build sequence
	for len(seq) < n {
		bestVal := -1.0
		bestIdx := -1
		for j := 1; j < n; j++ {
			if use[j] {
				continue
			}
			seq = append(seq, j)
			val := calc(seq)
			if val > bestVal {
				bestVal = val
				bestIdx = j
			}
			seq = seq[:len(seq)-1]
		}
		seq = append(seq, bestIdx)
		use[bestIdx] = true
	}
	result := calc(seq)
	fmt.Printf("%.20f\n", result)
}
