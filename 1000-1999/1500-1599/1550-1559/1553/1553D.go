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
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)

		i := len(s) - 1
		j := len(t) - 1
		for i >= 0 && j >= 0 {
			if s[i] == t[j] {
				i--
				j--
			} else {
				i -= 2
			}
		}

		if j < 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
