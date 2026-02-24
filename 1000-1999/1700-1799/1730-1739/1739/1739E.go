package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)

	var row1, row2 string
	fmt.Fscan(reader, &row1)
	fmt.Fscan(reader, &row2)

	dp := [2][2]int{
		{0, -1000000},
		{-1000000, -1000000},
	}

	for c := 0; c < n; c++ {
		var grid [2]int
		grid[0] = int(row1[c] - '0')
		grid[1] = int(row2[c] - '0')

		if c == 0 {
			grid[0] = 0
		}

		new_dp := [2][2]int{
			{-1000000, -1000000},
			{-1000000, -1000000},
		}

		for r := 0; r < 2; r++ {
			for f := 0; f < 2; f++ {
				if dp[r][f] < 0 {
					continue
				}

				val_r := grid[r]
				val_o := grid[1-r]
				if f == 1 {
					val_o = 0
				}

				if dp[r][f] > new_dp[r][0] {
					new_dp[r][0] = dp[r][f]
				}

				if val_r == 1 {
					if dp[r][f]+1 > new_dp[r][0] {
						new_dp[r][0] = dp[r][f] + 1
					}
				}

				if val_o == 1 {
					if dp[r][f]+1 > new_dp[1-r][1] {
						new_dp[1-r][1] = dp[r][f] + 1
					}
				}

				if val_r == 1 && val_o == 1 {
					if dp[r][f]+2 > new_dp[1-r][1] {
						new_dp[1-r][1] = dp[r][f] + 2
					}
				}
			}
		}
		dp = new_dp
	}

	ans := 0
	for r := 0; r < 2; r++ {
		for f := 0; f < 2; f++ {
			if dp[r][f] > ans {
				ans = dp[r][f]
			}
		}
	}

	fmt.Println(ans)
}