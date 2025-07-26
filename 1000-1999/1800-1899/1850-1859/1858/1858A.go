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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		if a > b {
			fmt.Fprintln(writer, "First")
		} else if b > a {
			fmt.Fprintln(writer, "Second")
		} else {
			if c%2 == 1 {
				fmt.Fprintln(writer, "First")
			} else {
				fmt.Fprintln(writer, "Second")
			}
		}
	}
}
