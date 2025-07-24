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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		ones := 0
		minBlock := 4
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '1' {
					ones++
				}
			}
		}
		for i := 0; i < n-1; i++ {
			for j := 0; j < m-1; j++ {
				cnt := 0
				if grid[i][j] == '1' {
					cnt++
				}
				if grid[i][j+1] == '1' {
					cnt++
				}
				if grid[i+1][j] == '1' {
					cnt++
				}
				if grid[i+1][j+1] == '1' {
					cnt++
				}
				if cnt < minBlock {
					minBlock = cnt
				}
			}
		}

		var ans int
		if ones == 0 {
			ans = 0
		} else if minBlock <= 2 {
			ans = ones
		} else if minBlock == 3 {
			ans = ones - 1
		} else {
			ans = ones - 2
		}
		fmt.Fprintln(out, ans)
	}
}
