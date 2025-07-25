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
		var a, b, n, s int64
		fmt.Fscan(reader, &a, &b, &n, &s)
		use := s / n
		if use > a {
			use = a
		}
		if s-use*n <= b {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
