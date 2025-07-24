package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	const mod int = 1000000007
	maxVal := n
	if m > maxVal {
		maxVal = m
	}
	dp := make([]int, maxVal+2)
	dp[1] = 2
	if maxVal >= 2 {
		dp[2] = 4
	}
	for i := 3; i <= maxVal; i++ {
		dp[i] = (dp[i-1] + dp[i-2]) % mod
	}

	ans := (dp[n] + dp[m] - 2) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
