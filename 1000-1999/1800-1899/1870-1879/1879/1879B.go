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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		sumA, sumB := int64(0), int64(0)
		minA, minB := a[0], b[0]
		for i := 0; i < n; i++ {
			sumA += a[i]
			sumB += b[i]
			if a[i] < minA {
				minA = a[i]
			}
			if b[i] < minB {
				minB = b[i]
			}
		}
		n64 := int64(n)
		costRows := sumA + n64*minB
		costCols := sumB + n64*minA
		if costRows < costCols {
			fmt.Fprintln(writer, costRows)
		} else {
			fmt.Fprintln(writer, costCols)
		}
	}
}
