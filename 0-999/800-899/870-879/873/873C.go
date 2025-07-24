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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &matrix[i][j])
		}
	}

	totalScore := 0
	totalRemove := 0

	for col := 0; col < m; col++ {
		prefix := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + matrix[i-1][col]
		}

		bestScore := 0
		bestRemove := 0
		for row := 1; row <= n; row++ {
			if matrix[row-1][col] == 1 {
				onesAbove := prefix[row-1]
				end := row + k - 1
				if end > n {
					end = n
				}
				score := prefix[end] - prefix[row-1]
				if score > bestScore || (score == bestScore && onesAbove < bestRemove) {
					bestScore = score
					bestRemove = onesAbove
				}
			}
		}

		totalScore += bestScore
		totalRemove += bestRemove
	}

	fmt.Fprintln(writer, totalScore, totalRemove)
}
