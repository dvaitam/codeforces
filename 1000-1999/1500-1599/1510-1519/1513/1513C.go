package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	const maxM = 200000
	dp := make([]int64, maxM+10+1)
	for i := 0; i <= 9; i++ {
		dp[i] = 1
	}
	for i := 10; i <= maxM+9; i++ {
		dp[i] = (dp[i-9] + dp[i-10]) % mod
	}

	for ; t > 0; t-- {
		var n string
		var m int
		fmt.Fscan(reader, &n, &m)
		var ans int64
		for _, ch := range n {
			d := int(ch - '0')
			ans = (ans + dp[d+m]) % mod
		}
		fmt.Fprintln(writer, ans)
	}
}
