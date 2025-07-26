package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		score := 0
		for i := 0; i < 10; i++ {
			var row string
			fmt.Fscan(in, &row)
			for j := 0; j < 10 && j < len(row); j++ {
				if row[j] == 'X' {
					val := min(min(i, 9-i), min(j, 9-j)) + 1
					score += val
				}
			}
		}
		fmt.Fprintln(out, score)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
