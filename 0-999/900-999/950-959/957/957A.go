package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	colors := []byte{'C', 'M', 'Y'}
	// dp[c] = number of ways (capped at 2) that last painted color is colors[c]
	dp := [3]int{}
	// initialize for first character
	for i, col := range colors {
		if s[0] == '?' || byte(s[0]) == col {
			dp[i] = 1
		}
	}

	for i := 1; i < n; i++ {
		newDp := [3]int{}
		for j, col := range colors {
			if s[i] != '?' && byte(s[i]) != col {
				continue
			}
			// sum over previous colors not equal to j
			for k := 0; k < 3; k++ {
				if k == j {
					continue
				}
				newDp[j] += dp[k]
				if newDp[j] > 2 {
					newDp[j] = 2
				}
			}
		}
		dp = newDp
	}

	total := 0
	for i := 0; i < 3; i++ {
		total += dp[i]
		if total > 2 {
			total = 2
			break
		}
	}
	if total >= 2 {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
