package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

var fib [55]int64
var dp [1005][55005]int64

func main() {
	fib[1], fib[2] = 1, 1
	for i := 3; i <= 30; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, x, m int
	fmt.Fscan(reader, &n, &x, &m)

	dp[0][0] = 1
	for i := 1; i <= x; i++ {
		for j := 1; j <= n; j++ {
			for l := fib[i]; l <= fib[i]*int64(j); l++ {
				dp[j][l] = (dp[j][l] + dp[j-1][l-fib[i]]) % mod
			}
		}
	}

	var ans int64
	limit := int(fib[x] * int64(n))
	for i := 0; i <= limit; i++ {
		t := i
		c := 0
		for j := 30; j >= 1; j-- {
			c += t / int(fib[j])
			t %= int(fib[j])
		}
		if c == m {
			ans = (ans + dp[n][i]) % mod
		}
	}

	fmt.Fprintln(writer, ans)
}
