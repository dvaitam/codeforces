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

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var n int
		var a rune
		var s string
		fmt.Fscan(reader, &n, &a, &s)

		count := 0
		lastPos := 0
		for i, ch := range s {
			if ch == a {
				count++
				lastPos = i + 1
			}
		}
		if count == n {
			fmt.Fprintln(writer, 0)
		} else if lastPos > n/2 || count == n-1 {
			fmt.Fprintln(writer, 1)
			fmt.Fprintln(writer, lastPos)
		} else {
			fmt.Fprintln(writer, 2)
			fmt.Fprintln(writer, n, n-1)
		}
	}
}
