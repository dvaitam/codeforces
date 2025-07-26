package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)

		dp := make([]int, 26)
		for i := 0; i < 26; i++ {
			dp[i] = math.MinInt32
		}
		dpNo := 0

		for _, ch := range s {
			idx := int(ch - 'a')
			newDp := make([]int, 26)
			for i := 0; i < 26; i++ {
				newDp[i] = dp[i]
			}
			newDpNo := dpNo

			if newDp[idx] < dpNo {
				newDp[idx] = dpNo
			}
			for c := 0; c < 26; c++ {
				if newDp[idx] < dp[c] {
					newDp[idx] = dp[c]
				}
			}
			if dp[idx] != math.MinInt32 && dp[idx]+2 > newDpNo {
				newDpNo = dp[idx] + 2
			}

			dp = newDp
			dpNo = newDpNo
		}

		result := len(s) - dpNo
		fmt.Fprintln(writer, result)
	}
}
