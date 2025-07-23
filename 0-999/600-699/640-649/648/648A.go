package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grid[i])
	}

	heights := make([]int, m)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if grid[i][j] == '*' {
				heights[j] = n - i
				break
			}
		}
	}

	maxUp, maxDown := 0, 0
	for j := 0; j < m-1; j++ {
		diff := heights[j+1] - heights[j]
		if diff > maxUp {
			maxUp = diff
		} else if -diff > maxDown {
			maxDown = -diff
		}
	}
	fmt.Println(maxUp, maxDown)
}
