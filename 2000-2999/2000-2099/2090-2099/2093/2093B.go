package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		sumDigits := 0
		for i := 0; i < n; i++ {
			sumDigits += int(s[i] - '0')
		}
		dp := make([]map[int]int, n+1)
		dp[0] = map[int]int{0: 0}
		for i := 0; i < n; i++ {
			dp[i+1] = make(map[int]int)
			digit := int(s[i] - '0')
			for sum, best := range dp[i] {
				if cur, ok := dp[i+1][sum]; !ok || cur > best+1 {
					dp[i+1][sum] = best + 1
				}
				dp[i+1][sum+digit] = best
			}
		}
		bestAns := n
		for sum := range dp[n] {
			remove := dp[n][sum]
			if sum > 0 && remove < bestAns {
				bestAns = remove
			}
		}
		fmt.Fprintln(out, bestAns)
	}
}
