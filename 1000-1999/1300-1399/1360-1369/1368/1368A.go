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
		var a, b, n int64
		fmt.Fscan(reader, &a, &b, &n)
		steps := 0
		for a <= n && b <= n {
			if a < b {
				a += b
			} else {
				b += a
			}
			steps++
		}
		fmt.Fprintln(writer, steps)
	}
}
