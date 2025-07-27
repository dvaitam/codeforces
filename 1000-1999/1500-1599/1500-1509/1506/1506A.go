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
		var n, m, x int64
		fmt.Fscan(reader, &n, &m, &x)
		x-- // zero-based index in column-major order
		row := x % n
		col := x / n
		ans := row*m + col + 1
		fmt.Fprintln(writer, ans)
	}
}
