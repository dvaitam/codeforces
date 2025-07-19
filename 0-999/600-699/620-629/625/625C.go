package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var n, m int
	if _, err := fmt.Fscan(os.Stdin, &n, &m); err != nil {
		return
	}
	// Initialize table
	ans := make([][]int, n)
	for i := range ans {
		ans[i] = make([]int, n)
	}
	times := 0
	sum := 0
	// Fill first m-1 columns by column
	for c := 0; c < m-1; c++ {
		for r := 0; r < n; r++ {
			times++
			ans[r][c] = times
		}
	}
	// Fill columns m-1 to n-1 by row, accumulate sum at column m-1
	for r := 0; r < n; r++ {
		for c := m - 1; c < n; c++ {
			times++
			ans[r][c] = times
			if c == m-1 {
				sum += times
			}
		}
	}
	// Output
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	w.WriteString(strconv.Itoa(sum))
	w.WriteByte('\n')
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			w.WriteString(strconv.Itoa(ans[r][c]))
			if c < n {
				w.WriteByte(' ')
			}
		}
		w.WriteByte('\n')
	}
}
