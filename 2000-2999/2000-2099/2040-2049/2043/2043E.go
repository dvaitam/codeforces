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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		A := make([][]int, n)
		B := make([][]int, n)
		for i := 0; i < n; i++ {
			A[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &A[i][j])
			}
		}
		for i := 0; i < n; i++ {
			B[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &B[i][j])
			}
		}

		possible := true
		for bit := 0; bit < 30 && possible; bit++ {
			if !checkBit(A, B, n, m, bit) {
				possible = false
			}
		}

		if possible {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}

func checkBit(A, B [][]int, n, m, bit int) bool {
	rowZero := make([]bool, n)
	colOr := make([]bool, m)
	rowQ := make([]int, 0)
	colQ := make([]int, 0)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			aBit := (A[i][j] >> bit) & 1
			bBit := (B[i][j] >> bit) & 1
			if aBit == 1 && bBit == 0 && !rowZero[i] {
				rowZero[i] = true
				rowQ = append(rowQ, i)
			}
			if aBit == 0 && bBit == 1 && !colOr[j] {
				colOr[j] = true
				colQ = append(colQ, j)
			}
		}
	}

	headR, headC := 0, 0
	for headR < len(rowQ) || headC < len(colQ) {
		for headR < len(rowQ) {
			i := rowQ[headR]
			headR++
			for j := 0; j < m; j++ {
				if colOr[j] {
					continue
				}
				if ((B[i][j] >> bit) & 1) == 1 {
					colOr[j] = true
					colQ = append(colQ, j)
				}
			}
		}
		for headC < len(colQ) {
			j := colQ[headC]
			headC++
			for i := 0; i < n; i++ {
				if rowZero[i] {
					continue
				}
				if ((B[i][j] >> bit) & 1) == 0 {
					rowZero[i] = true
					rowQ = append(rowQ, i)
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			bBit := (B[i][j] >> bit) & 1
			if rowZero[i] {
				if !colOr[j] && bBit == 1 {
					return false
				}
			} else {
				if colOr[j] {
					if bBit == 0 {
						return false
					}
				} else {
					aBit := (A[i][j] >> bit) & 1
					if bBit != aBit {
						return false
					}
				}
			}
		}
	}

	rows := make([]int, 0)
	for i := 0; i < n; i++ {
		if rowZero[i] {
			rows = append(rows, i)
		}
	}
	cols := make([]int, 0)
	for j := 0; j < m; j++ {
		if colOr[j] {
			cols = append(cols, j)
		}
	}

	if len(cols) == 0 || len(rows) <= 1 {
		return true
	}
	for x := 0; x < len(rows); x++ {
		for y := x + 1; y < len(rows); y++ {
			i := rows[x]
			k := rows[y]
			le := true
			ge := true
			for _, j := range cols {
				bi := (B[i][j] >> bit) & 1
				bk := (B[k][j] >> bit) & 1
				if bi < bk {
					ge = false
				}
				if bi > bk {
					le = false
				}
				if !le && !ge {
					return false
				}
			}
		}
	}
	return true
}
