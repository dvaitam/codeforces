package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	rowMask := make([]uint64, n)
	colMask := make([]uint64, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				rowMask[i] |= 1 << uint(j)
				colMask[j] |= 1 << uint(i)
			}
		}
	}

	// check rows: any two rows that share a column must have identical masks
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rowMask[i]&rowMask[j] != 0 && rowMask[i] != rowMask[j] {
				fmt.Println("No")
				return
			}
		}
	}

	// check columns: any two columns that share a row must have identical masks
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			if colMask[i]&colMask[j] != 0 && colMask[i] != colMask[j] {
				fmt.Println("No")
				return
			}
		}
	}

	fmt.Println("Yes")
}
