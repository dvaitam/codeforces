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

	var n, m, x, y int
	fmt.Fscan(reader, &n, &m, &x, &y)

	// build initial column order: start at y, then all other columns
	v := make([]int, 0, m)
	v = append(v, y)
	for i := 1; i <= m; i++ {
		if i != y {
			v = append(v, i)
		}
	}

	// traverse rows from x to n, then 1 to x-1
	for i := x; i <= n; i++ {
		for _, col := range v {
			fmt.Fprintln(writer, i, col)
		}
		// reverse columns order for next row
		for l, r := 0, len(v)-1; l < r; l, r = l+1, r-1 {
			v[l], v[r] = v[r], v[l]
		}
	}
	for i := 1; i < x; i++ {
		for _, col := range v {
			fmt.Fprintln(writer, i, col)
		}
		for l, r := 0, len(v)-1; l < r; l, r = l+1, r-1 {
			v[l], v[r] = v[r], v[l]
		}
	}
}
