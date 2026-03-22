package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string

	_, err := fmt.Fscan(reader, &n, &s)
	if err != nil {
		return
	}

	if n%2 != 0 {
		fmt.Println(0)
		return
	}

	dp := make([]uint32, n/2+2)
	dp[0] = 1

	for i := 1; i <= n; i++ {
		isQ := (s[i-1] == '?')
		minJ := i % 2
		maxJ := i
		if n-i < maxJ {
			maxJ = n - i
		}

		j := minJ
		if j == 0 {
			if isQ {
				dp[0] = dp[1]
			} else {
				dp[0] = 0
			}
			j += 2
		}

		if isQ {
			for ; j <= maxJ; j += 2 {
				dp[j] = dp[j-1]*25 + dp[j+1]
			}
		} else {
			for ; j <= maxJ; j += 2 {
				dp[j] = dp[j-1]
			}
		}
	}

	fmt.Println(dp[0])
}
