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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &mat[i][j])
		}
	}

	arr := make([]int, n)
	for v := 1; v <= n; v++ {
		for idx := 0; idx < n; idx++ {
			if arr[idx] != 0 {
				continue
			}
			c := 0
			for j := 0; j < n; j++ {
				if mat[idx][j] == v {
					c++
				}
			}
			if c == n-v {
				arr[idx] = v
				break
			}
		}
	}

	for i, v := range arr {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
