package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	fmt.Fscan(in, &k)
	n := 1 << k
	mat := [][]byte{[]byte{'+'}}
	size := 1
	for size < n {
		newSize := size * 2
		newMat := make([][]byte, newSize)
		for i := range newMat {
			newMat[i] = make([]byte, newSize)
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				ch := mat[i][j]
				newMat[i][j] = ch
				newMat[i][j+size] = ch
				newMat[i+size][j] = ch
				if ch == '+' {
					newMat[i+size][j+size] = '*'
				} else {
					newMat[i+size][j+size] = '+'
				}
			}
		}
		mat = newMat
		size = newSize
	}
	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		out.Write(mat[i])
		out.WriteByte('\n')
	}
	out.Flush()
}
