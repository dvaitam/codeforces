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
	var s string
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}

	u, d, l, r := 0, 0, 0, 0
	for _, ch := range s {
		switch ch {
		case 'U':
			u++
		case 'D':
			d++
		case 'L':
			l++
		case 'R':
			r++
		}
	}

	// Maximum commands executed correctly are pairs of opposite directions.
	res := 2*min(u, d) + 2*min(l, r)
	fmt.Fprintln(writer, res)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
