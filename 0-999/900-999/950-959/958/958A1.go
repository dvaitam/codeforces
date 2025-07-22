package main

import (
	"bufio"
	"fmt"
	"os"
)

func transform(mat [][]byte, op int) [][]byte {
	n := len(mat)
	res := make([][]byte, n)
	for i := range res {
		res[i] = make([]byte, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			switch op {
			case 0: // identity
				res[i][j] = mat[i][j]
			case 1: // rotate 90 clockwise
				res[i][j] = mat[n-1-j][i]
			case 2: // rotate 180
				res[i][j] = mat[n-1-i][n-1-j]
			case 3: // rotate 270 clockwise
				res[i][j] = mat[j][n-1-i]
			case 4: // vertical flip
				res[i][j] = mat[i][n-1-j]
			case 5: // horizontal flip
				res[i][j] = mat[n-1-i][j]
			case 6: // main diagonal flip
				res[i][j] = mat[j][i]
			case 7: // anti diagonal flip
				res[i][j] = mat[n-1-j][n-1-i]
			}
		}
	}
	return res
}

func equal(a, b [][]byte) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([][]byte, n)
	b := make([][]byte, n)
	var line string
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &line)
		a[i] = []byte(line)
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &line)
		b[i] = []byte(line)
	}

	for op := 0; op < 8; op++ {
		t := transform(a, op)
		if equal(t, b) {
			fmt.Fprintln(writer, "Yes")
			return
		}
	}
	fmt.Fprintln(writer, "No")
}
