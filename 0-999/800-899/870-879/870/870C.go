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
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int64
		fmt.Fscan(reader, &n)
		switch {
		case n < 4:
			fmt.Fprintln(writer, -1)
		case n%4 == 0:
			fmt.Fprintln(writer, n/4)
		case n%4 == 1:
			if n >= 9 {
				fmt.Fprintln(writer, (n-9)/4+1)
			} else {
				fmt.Fprintln(writer, -1)
			}
		case n%4 == 2:
			if n >= 6 {
				fmt.Fprintln(writer, (n-6)/4+1)
			} else {
				fmt.Fprintln(writer, -1)
			}
		default: // n%4 == 3
			if n >= 15 {
				fmt.Fprintln(writer, (n-15)/4+2)
			} else {
				fmt.Fprintln(writer, -1)
			}
		}
	}
}
