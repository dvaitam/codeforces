package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	if n <= 0 {
		return
	}

	const maxA = 1000000
	presence := make([]bool, maxA+1)
	numbers := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &numbers[i])
		if numbers[i] > maxVal {
			maxVal = numbers[i]
		}
		presence[numbers[i]] = true
	}

	dp := make([]int, maxVal+1)
	ans := 1
	for x := 1; x <= maxVal; x++ {
		if !presence[x] {
			continue
		}
		if dp[x] == 0 {
			dp[x] = 1
		}
		if dp[x] > ans {
			ans = dp[x]
		}
		for m := x * 2; m <= maxVal; m += x {
			if presence[m] {
				if dp[m] < dp[x]+1 {
					dp[m] = dp[x] + 1
				}
			}
		}
	}

	fmt.Println(ans)
}
