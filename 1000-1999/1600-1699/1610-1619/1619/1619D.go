package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(matrix [][]int, m, n int, x int) bool {
	friendOk := make([]bool, n)
	hasDouble := false
	for i := 0; i < m; i++ {
		cnt := 0
		row := matrix[i]
		for j := 0; j < n; j++ {
			if row[j] >= x {
				if !friendOk[j] {
					friendOk[j] = true
				}
				cnt++
			}
		}
		if cnt >= 2 {
			hasDouble = true
		}
	}
	if !hasDouble {
		return false
	}
	for j := 0; j < n; j++ {
		if !friendOk[j] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m, n int
		fmt.Fscan(reader, &m, &n)
		matrix := make([][]int, m)
		for i := 0; i < m; i++ {
			matrix[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &matrix[i][j])
			}
		}
		low := 1
		high := 1000000000
		for low < high {
			mid := (low + high + 1) / 2
			if check(matrix, m, n, mid) {
				low = mid
			} else {
				high = mid - 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}
