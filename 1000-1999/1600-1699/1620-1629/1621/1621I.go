package main

import (
	"bufio"
	"fmt"
	"os"
)

func lexSmaller(a []int, i1 int, b []int, i2 int, length int) bool {
	for k := 0; k < length; k++ {
		if a[i1+k] < b[i2+k] {
			return true
		} else if a[i1+k] > b[i2+k] {
			return false
		}
	}
	return false
}

func op(arr []int) []int {
	n := len(arr)
	D := make([]int, n)
	copy(D, arr)
	for i := 1; i <= n; i++ {
		bestIdx := 0
		for s := 1; s <= n-i; s++ {
			if lexSmaller(D, s, D, bestIdx, i) {
				bestIdx = s
			}
		}
		copy(D[n-i:], D[bestIdx:bestIdx+i])
	}
	return D
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
	}
	var q int
	fmt.Fscan(reader, &q)

	// precompute B arrays
	B := make([][]int, n+1)
	B[0] = make([]int, n)
	copy(B[0], A)
	for i := 1; i <= n; i++ {
		B[i] = op(B[i-1])
	}

	for ; q > 0; q-- {
		var i, j int
		fmt.Fscan(reader, &i, &j)
		if i > n {
			i = n
		}
		if j > n {
			j = n
		}
		fmt.Fprintln(writer, B[i][j-1])
	}
}
