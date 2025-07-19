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
		var n int
		fmt.Fscan(reader, &n)
		for i := 0; i < n; i++ {
			writer.WriteByte('1')
			if i+1 < n {
				writer.WriteByte(' ')
			}
		}
		writer.WriteByte('\n')
	}
}
