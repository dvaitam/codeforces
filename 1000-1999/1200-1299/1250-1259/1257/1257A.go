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
		var n, x, a, b int
		fmt.Fscan(reader, &n, &x, &a, &b)
		if a > b {
			a, b = b, a
		}
		dist := b - a
		if dist+x > n-1 {
			fmt.Fprintln(writer, n-1)
		} else {
			fmt.Fprintln(writer, dist+x)
		}
	}
}
