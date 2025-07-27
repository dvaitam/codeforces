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
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		// prefix sums of ones for each row
		prefix := make([][]int, n)
		for i := 0; i < n; i++ {
			prefix[i] = make([]int, m+1)
			for j := 0; j < m; j++ {
				prefix[i][j+1] = prefix[i][j]
				if grid[i][j] == '1' {
					prefix[i][j+1]++
				}
			}
		}

		ans := n * m // upper bound
		if n >= 5 && m >= 4 {
			for top := 0; top+4 < n; top++ {
				bottom := top + 4
				// precompute mid column prefix sums and ones counts
				midOnes := make([]int, m)
				midPrefix := make([]int, m+1)
				for j := 0; j < m; j++ {
					cnt := 0
					for r := top + 1; r <= top+3; r++ {
						if grid[r][j] == '1' {
							cnt++
						}
					}
					midOnes[j] = cnt
					midPrefix[j+1] = midPrefix[j] + cnt
				}
				for left := 0; left+3 < m; left++ {
					for right := left + 3; right < m; right++ {
						width := right - left - 1
						// top and bottom horizontal edges cost
						topOnes := prefix[top][right] - prefix[top][left+1]
						bottomOnes := prefix[bottom][right] - prefix[bottom][left+1]
						costTop := width - topOnes
						costBottom := width - bottomOnes
						// interior zero cost
						interiorOnes := midPrefix[right] - midPrefix[left+1]
						// vertical edges cost
						verticalCost := (3 - midOnes[left]) + (3 - midOnes[right])
						cost := costTop + costBottom + interiorOnes + verticalCost
						if cost < ans {
							ans = cost
						}
					}
				}
			}
		} else {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}
