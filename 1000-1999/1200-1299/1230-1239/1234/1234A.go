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
		var n int
		fmt.Fscan(reader, &n)
		sum := int64(0)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			sum += x
		}
		price := (sum + int64(n) - 1) / int64(n)
		fmt.Fprintln(writer, price)
	}
}
