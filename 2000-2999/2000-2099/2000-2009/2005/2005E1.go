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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, n, m int
		fmt.Fscan(in, &l, &n, &m)

		a := make([]int, l)
		for i := 0; i < l; i++ {
			fmt.Fscan(in, &a[i])
		}

		b := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			b[i] = make([]int, m+1)
			for j := 1; j <= m; j++ {
				fmt.Fscan(in, &b[i][j])
			}
		}

		dpNext := makeBoolMatrix(n+2, m+2)

		for idx := l - 1; idx >= 0; idx-- {
			val := a[idx]
			diff := makeIntMatrix(n+2, m+2)

			for r := 1; r <= n; r++ {
				row := b[r]
				dpRow := dpNext[r+1]
				for c := 1; c <= m; c++ {
					if row[c] == val && !dpRow[c+1] {
						diff[1][1]++
						diff[r+1][1]--
						diff[1][c+1]--
						diff[r+1][c+1]++
					}
				}
			}

			dpCur := makeBoolMatrix(n+2, m+2)
			for r := 1; r <= n+1; r++ {
				for c := 1; c <= m+1; c++ {
					diff[r][c] += diff[r-1][c] + diff[r][c-1] - diff[r-1][c-1]
					if diff[r][c] > 0 {
						dpCur[r][c] = true
					}
				}
			}

			dpNext = dpCur
		}

		if dpNext[1][1] {
			fmt.Fprintln(out, "T")
		} else {
			fmt.Fprintln(out, "N")
		}
	}
}

func makeBoolMatrix(n, m int) [][]bool {
	mat := make([][]bool, n)
	for i := range mat {
		mat[i] = make([]bool, m)
	}
	return mat
}

func makeIntMatrix(n, m int) [][]int {
	mat := make([][]int, n)
	for i := range mat {
		mat[i] = make([]int, m)
	}
	return mat
}
