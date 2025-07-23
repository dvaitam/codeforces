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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	pow2 := make([]int64, 51)
	pow2[0] = 1
	for i := 1; i <= 50; i++ {
		pow2[i] = pow2[i-1] * 2
	}

	ans := int64(0)
	// count subsets for rows
	for i := 0; i < n; i++ {
		cnt0, cnt1 := 0, 0
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 {
				cnt0++
			} else {
				cnt1++
			}
		}
		ans += pow2[cnt0] - 1
		ans += pow2[cnt1] - 1
	}
	// count subsets for columns
	for j := 0; j < m; j++ {
		cnt0, cnt1 := 0, 0
		for i := 0; i < n; i++ {
			if grid[i][j] == 0 {
				cnt0++
			} else {
				cnt1++
			}
		}
		ans += pow2[cnt0] - 1
		ans += pow2[cnt1] - 1
	}

	ans -= int64(n * m) // single cells counted twice
	fmt.Fprintln(out, ans)
}
