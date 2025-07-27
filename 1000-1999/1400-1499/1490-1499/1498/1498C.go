package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	tests := make([][2]int, t)
	maxN, maxK := 0, 0
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &tests[i][0], &tests[i][1])
		n, k := tests[i][0], tests[i][1]
		if n > maxN {
			maxN = n
		}
		if k > maxK {
			maxK = k
		}
	}

	dp := make([][]int64, maxN+1)
	for i := range dp {
		dp[i] = make([]int64, maxK+1)
	}
	for n := 0; n <= maxN; n++ {
		dp[n][1] = 1
	}
	for k := 1; k <= maxK; k++ {
		dp[0][k] = 1
	}
	for k := 2; k <= maxK; k++ {
		prefix := dp[0][k-1] % MOD
		for n := 1; n <= maxN; n++ {
			dp[n][k] = (1 + prefix) % MOD
			prefix = (prefix + dp[n][k-1]) % MOD
		}
	}

	for _, test := range tests {
		n, k := test[0], test[1]
		fmt.Fprintln(writer, dp[n][k]%MOD)
	}
}
