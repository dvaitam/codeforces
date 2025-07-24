package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var a, b int64
		fmt.Fscan(reader, &n, &a, &b)
		var s string
		fmt.Fscan(reader, &s)
		dp0 := make([]int64, n+1)
		dp1 := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			dp0[i] = inf
			dp1[i] = inf
		}
		dp0[0] = b
		for i := 1; i <= n; i++ {
			if s[i-1] == '1' {
				dp0[i] = inf
				dp1[i] = dp1[i-1] + a + 2*b
			} else {
				dp0[i] = min(dp0[i-1]+a+b, dp1[i-1]+2*a+b)
				dp1[i] = min(dp1[i-1]+a+2*b, dp0[i-1]+2*a+2*b)
			}
		}
		fmt.Fprintln(writer, dp0[n])
	}
}
