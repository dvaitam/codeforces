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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	cmds := make([]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		cmds[i] = s[0]
	}

	const mod int = 1e9 + 7
	dp := make([]int, n+2)
	dp[0] = 1

	for _, c := range cmds {
		if c == 'f' {
			// next line must be indented one level deeper
			newDP := make([]int, n+2)
			for j := 0; j <= n; j++ {
				if dp[j] != 0 {
					newDP[j+1] = (newDP[j+1] + dp[j]) % mod
				}
			}
			dp = newDP
		} else { // simple statement
			for j := n; j >= 0; j-- {
				dp[j] %= mod
				if j < n {
					dp[j] = (dp[j] + dp[j+1]) % mod
				}
			}
		}
	}

	fmt.Fprintln(writer, dp[0]%mod)
}
