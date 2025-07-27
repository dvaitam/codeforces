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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	const maxVal = 200000
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make([]int, maxVal+1)
		maxA := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x > maxA {
				maxA = x
			}
			freq[x]++
		}
		dp := make([]int, maxA+1)
		for i := 1; i <= maxA; i++ {
			dp[i] = freq[i]
		}
		for i := 1; i <= maxA; i++ {
			if freq[i] == 0 {
				continue
			}
			for j := i * 2; j <= maxA; j += i {
				if freq[j] > 0 && dp[j] < dp[i]+freq[j] {
					dp[j] = dp[i] + freq[j]
				}
			}
		}
		best := 0
		for i := 1; i <= maxA; i++ {
			if dp[i] > best {
				best = dp[i]
			}
		}
		fmt.Fprintln(writer, n-best)
	}
}
