package main

import (
	"bufio"
	"fmt"
	"os"
)

func processDiag(mat [][]int64, n int, startRow, startCol int) int64 {
	minVal := int64(1<<62 - 1)
	i, j := startRow, startCol
	for i < n && j < n {
		if mat[i][j] < minVal {
			minVal = mat[i][j]
		}
		i++
		j++
	}
	if minVal < 0 {
		return -minVal
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		mat := make([][]int64, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &mat[i][j])
			}
		}
		var ans int64
		for col := 0; col < n; col++ {
			ans += processDiag(mat, n, 0, col)
		}
		for row := 1; row < n; row++ {
			ans += processDiag(mat, n, row, 0)
		}
		fmt.Fprintln(out, ans)
	}
}
