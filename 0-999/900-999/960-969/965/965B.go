package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	cnt := make([][]int, n)
	for i := range cnt {
		cnt[i] = make([]int, n)
	}

	// horizontal placements
	for i := 0; i < n; i++ {
		for j := 0; j+k <= n; j++ {
			ok := true
			for t := 0; t < k; t++ {
				if grid[i][j+t] == '#' {
					ok = false
					break
				}
			}
			if ok {
				for t := 0; t < k; t++ {
					cnt[i][j+t]++
				}
			}
		}
	}

	// vertical placements
	for j := 0; j < n; j++ {
		for i := 0; i+k <= n; i++ {
			ok := true
			for t := 0; t < k; t++ {
				if grid[i+t][j] == '#' {
					ok = false
					break
				}
			}
			if ok {
				for t := 0; t < k; t++ {
					cnt[i+t][j]++
				}
			}
		}
	}

	bestRow, bestCol, bestVal := 0, 0, -1
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if cnt[i][j] > bestVal {
				bestVal = cnt[i][j]
				bestRow = i
				bestCol = j
			}
		}
	}

	fmt.Fprintf(writer, "%d %d\n", bestRow+1, bestCol+1)
}
