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
		var sumWidth int64
		var maxHeight int64
		for i := 0; i < n; i++ {
			var a, b int64
			fmt.Fscan(reader, &a, &b)
			if a > b {
				a, b = b, a
			}
			sumWidth += a
			if b > maxHeight {
				maxHeight = b
			}
		}
		perimeter := 2 * (sumWidth + maxHeight)
		fmt.Fprintln(writer, perimeter)
	}
}
