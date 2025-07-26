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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		if n%2 == 1 || n < 4 {
			fmt.Fprintln(writer, -1)
			continue
		}
		min := (n + 5) / 6
		max := n / 4
		fmt.Fprintln(writer, min, max)
	}
}
