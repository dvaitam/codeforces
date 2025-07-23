package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	for i := 0; i < n && k > 0; i++ {
		if k > 0 {
			matrix[i][i] = 1
			k--
		}
		for j := i + 1; j < n && k >= 2; j++ {
			matrix[i][j] = 1
			matrix[j][i] = 1
			k -= 2
		}
	}

	if k > 0 {
		fmt.Println(-1)
		return
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(matrix[i][j])
		}
		fmt.Println()
	}
}
