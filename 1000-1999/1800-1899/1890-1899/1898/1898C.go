package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	var T int
	if _, err := fmt.Fscan(br, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var N, M, K int
		fmt.Fscan(br, &N, &M, &K)
		rows := 2*N - 1
		cols := 2*M - 1
		// initialize grid
		C := make([][]byte, rows)
		for i := 0; i < rows; i++ {
			C[i] = make([]byte, cols)
			for j := 0; j < cols; j++ {
				C[i][j] = '#'
			}
		}
		// check feasibility
		minMoves := N + M - 2
		if K < minMoves || (K%2) != (minMoves%2) {
			bw.WriteString("NO\n")
			continue
		}
		maxRow := 2*N - 2
		maxCol := 2*M - 2
		// fill paths
		var now int
		if K%4 == minMoves%4 {
			now = 0
			for j := 1; j < maxCol; j += 2 {
				if now == 0 {
					C[0][j] = 'R'
				} else {
					C[0][j] = 'B'
				}
				now ^= 1
			}
			for i := 1; i < maxRow; i += 2 {
				if now == 0 {
					C[i][maxCol] = 'R'
				} else {
					C[i][maxCol] = 'B'
				}
				now ^= 1
			}
			if now == 0 {
				C[maxRow-1][maxCol] = 'B'
				C[maxRow-1][maxCol-2] = 'B'
				C[maxRow-2][maxCol-1] = 'R'
				C[maxRow][maxCol-1] = 'R'
			} else {
				C[maxRow-1][maxCol] = 'R'
				C[maxRow-1][maxCol-2] = 'R'
				C[maxRow-2][maxCol-1] = 'B'
				C[maxRow][maxCol-1] = 'B'
			}
		} else {
			now = 0
			for j := 1; j < maxCol; j += 2 {
				if now == 0 {
					C[0][j] = 'R'
				} else {
					C[0][j] = 'B'
				}
				now ^= 1
			}
			for i := 1; i < maxRow; i += 2 {
				if now == 0 {
					C[i][maxCol] = 'R'
				} else {
					C[i][maxCol] = 'B'
				}
				now ^= 1
			}
			if now == 1 {
				C[maxRow-1][maxCol] = 'B'
				C[maxRow-1][maxCol-2] = 'B'
				C[maxRow-2][maxCol-1] = 'R'
				C[maxRow][maxCol-1] = 'R'
			} else {
				C[maxRow-1][maxCol] = 'R'
				C[maxRow-1][maxCol-2] = 'R'
				C[maxRow-2][maxCol-1] = 'B'
				C[maxRow][maxCol-1] = 'B'
			}
		}
		// output
		bw.WriteString("YES\n")
		// first pattern rows
		for i := 0; i <= maxRow; i += 2 {
			for j := 1; j <= maxCol; j += 2 {
				if C[i][j] == 'B' {
					bw.WriteString("B ")
				} else {
					bw.WriteString("R ")
				}
			}
			bw.WriteByte('\n')
		}
		// second pattern rows
		for i := 1; i <= maxRow; i += 2 {
			for j := 0; j <= maxCol; j += 2 {
				if C[i][j] == 'B' {
					bw.WriteString("B ")
				} else {
					bw.WriteString("R ")
				}
			}
			bw.WriteByte('\n')
		}
	}
}
