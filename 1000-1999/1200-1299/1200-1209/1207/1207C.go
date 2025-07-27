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

		dp0, dp1 := b, inf
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				dp0 = inf
				dp1 = dp1 + a + 2*b
			} else {
				ndp0 := min(dp0+a+b, dp1+2*a+b)
				ndp1 := min(dp1+a+2*b, dp0+2*a+2*b)
				dp0, dp1 = ndp0, ndp1
			}
		}
		fmt.Fprintln(writer, dp0)
	}
}
