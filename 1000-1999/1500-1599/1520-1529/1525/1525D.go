package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	ones := make([]int, 0, n/2)
	zeros := make([]int, 0, n)
	for i, v := range arr {
		if v == 1 {
			ones = append(ones, i)
		} else {
			zeros = append(zeros, i)
		}
	}

	m := len(ones)
	if m == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	const inf = int(1e9)
	dp := make([]int, m+1)
	for i := 1; i <= m; i++ {
		dp[i] = inf
	}

	for _, pos0 := range zeros {
		for i := m; i >= 1; i-- {
			d := abs(ones[i-1] - pos0)
			if dp[i-1]+d < dp[i] {
				dp[i] = dp[i-1] + d
			}
		}
	}

	fmt.Fprintln(writer, dp[m])
}
