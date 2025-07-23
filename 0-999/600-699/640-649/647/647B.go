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

	var n, m int
	fmt.Fscan(in, &n, &m)

	minRow, maxRow := n, -1
	minCol, maxCol := m, -1

	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		for j, c := range line {
			if c == '*' {
				if i < minRow {
					minRow = i
				}
				if i > maxRow {
					maxRow = i
				}
				if j < minCol {
					minCol = j
				}
				if j > maxCol {
					maxCol = j
				}
			}
		}
	}

	height := maxRow - minRow + 1
	width := maxCol - minCol + 1

	if height > width {
		fmt.Fprintln(out, height)
	} else {
		fmt.Fprintln(out, width)
	}
}
