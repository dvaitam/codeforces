package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program simulates falling stones for problemG.txt in contest 1669.
// Stones ('*') fall vertically until hitting the grid bottom, an obstacle ('o'),
// or another stone.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			grid[i] = []byte(s)
		}

		result := make([][]byte, n)
		for i := range result {
			result[i] = make([]byte, m)
			for j := range result[i] {
				result[i][j] = '.'
			}
		}

		for c := 0; c < m; c++ {
			pos := n - 1
			for r := n - 1; r >= 0; r-- {
				switch grid[r][c] {
				case 'o':
					result[r][c] = 'o'
					pos = r - 1
				case '*':
					result[pos][c] = '*'
					pos--
				}
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprintln(writer, string(result[i]))
		}
	}
}
