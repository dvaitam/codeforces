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
	fmt.Fscan(in, &n, &m)

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	// reachable from start (1,1)
	fromStart := make([][]bool, n)
	for i := range fromStart {
		fromStart[i] = make([]bool, m)
	}
	if grid[0][0] != '#' {
		fromStart[0][0] = true
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				continue
			}
			if i == 0 && j == 0 {
				continue
			}
			if (i > 0 && fromStart[i-1][j]) || (j > 0 && fromStart[i][j-1]) {
				fromStart[i][j] = true
			}
		}
	}

	if !fromStart[n-1][m-1] {
		fmt.Fprintln(out, 0)
		return
	}

	// reachable to goal (n,m)
	toGoal := make([][]bool, n)
	for i := range toGoal {
		toGoal[i] = make([]bool, m)
	}
	if grid[n-1][m-1] != '#' {
		toGoal[n-1][m-1] = true
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if grid[i][j] == '#' {
				continue
			}
			if i == n-1 && j == m-1 {
				continue
			}
			if (i+1 < n && toGoal[i+1][j]) || (j+1 < m && toGoal[i][j+1]) {
				toGoal[i][j] = true
			}
		}
	}

	counts := make([]int, n+m+3)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if fromStart[i][j] && toGoal[i][j] {
				if (i == 0 && j == 0) || (i == n-1 && j == m-1) {
					continue
				}
				counts[i+j+2]++ // step index r+c where r,c are 1-based
			}
		}
	}

	for s := 3; s <= n+m-1; s++ {
		if counts[s] == 1 {
			fmt.Fprintln(out, 1)
			return
		}
	}
	fmt.Fprintln(out, 2)
}
