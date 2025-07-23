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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	rowTime := make([]int, n+1)
	rowColor := make([]int, n+1)
	colTime := make([]int, m+1)
	colColor := make([]int, m+1)

	for i := 1; i <= k; i++ {
		var t, idx int
		var color int
		fmt.Fscan(reader, &t, &idx, &color)
		if t == 1 {
			rowTime[idx] = i
			rowColor[idx] = color
		} else {
			colTime[idx] = i
			colColor[idx] = color
		}
	}

	for r := 1; r <= n; r++ {
		for c := 1; c <= m; c++ {
			val := 0
			if rowTime[r] > colTime[c] {
				val = rowColor[r]
			} else {
				val = colColor[c]
			}
			if c > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, val)
		}
		if r < n {
			fmt.Fprintln(writer)
		}
	}
}
