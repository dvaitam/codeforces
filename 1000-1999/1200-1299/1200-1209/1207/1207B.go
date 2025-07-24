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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	A := make([][]int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &A[i][j])
		}
	}

	B := make([][]int, n)
	for i := range B {
		B[i] = make([]int, m)
	}

	type op struct{ x, y int }
	var ops []op

	for i := 0; i+1 < n; i++ {
		for j := 0; j+1 < m; j++ {
			if A[i][j] == 1 && A[i][j+1] == 1 && A[i+1][j] == 1 && A[i+1][j+1] == 1 {
				ops = append(ops, op{i + 1, j + 1})
				B[i][j] = 1
				B[i][j+1] = 1
				B[i+1][j] = 1
				B[i+1][j+1] = 1
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] != B[i][j] {
				fmt.Fprintln(writer, -1)
				return
			}
		}
	}

	fmt.Fprintln(writer, len(ops))
	for _, o := range ops {
		fmt.Fprintln(writer, o.x, o.y)
	}
}
