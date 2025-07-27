package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		row1 := make([]int64, m)
		row2 := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &row1[i])
		}
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &row2[i])
		}
		prefix1 := make([]int64, m+1)
		prefix2 := make([]int64, m+1)
		for i := 0; i < m; i++ {
			prefix1[i+1] = prefix1[i] + row1[i]
			prefix2[i+1] = prefix2[i] + row2[i]
		}
		ans := int64(1<<63 - 1)
		for k := 0; k < m; k++ {
			top := prefix1[m] - prefix1[k+1]
			bottom := prefix2[k]
			best := top
			if bottom > best {
				best = bottom
			}
			ans = min(ans, best)
		}
		fmt.Fprintln(writer, ans)
	}
}
