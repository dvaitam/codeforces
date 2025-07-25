package main

import (
	"bufio"
	"fmt"
	"os"
)

func countSegments(s string) int {
	count := 0
	prevStar := false
	for i := 0; i < len(s); i++ {
		if s[i] == '*' {
			if !prevStar {
				count++
				prevStar = true
			}
		} else {
			prevStar = false
		}
	}
	return count
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
		var n int
		fmt.Fscan(reader, &n)
		var row1, row2 string
		fmt.Fscan(reader, &row1)
		fmt.Fscan(reader, &row2)

		segs := countSegments(row1) + countSegments(row2)
		ans := n - segs
		fmt.Fprintln(writer, ans)
	}
}
