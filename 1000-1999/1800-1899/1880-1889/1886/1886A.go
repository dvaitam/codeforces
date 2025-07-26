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
		var n int
		fmt.Fscan(reader, &n)
		if n < 7 || n == 9 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if n%3 == 0 {
			fmt.Fprintln(writer, "YES")
			fmt.Fprintln(writer, 1, 4, n-5)
		} else {
			fmt.Fprintln(writer, "YES")
			fmt.Fprintln(writer, 1, 2, n-3)
		}
	}
}
