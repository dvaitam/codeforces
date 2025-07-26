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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &rows[i])
		}

		diag := make([]int, n)
		total := 0
		for i := 0; i < n; i++ {
			row := rows[i]
			for j := 0; j < n; j++ {
				if row[j] == '1' {
					total++
					diff := (i - j + n) % n
					diag[diff]++
				}
			}
		}

		best := 0
		for i := 0; i < n; i++ {
			if diag[i] > best {
				best = diag[i]
			}
		}

		result := n + total - 2*best
		fmt.Fprintln(writer, result)
	}
}
