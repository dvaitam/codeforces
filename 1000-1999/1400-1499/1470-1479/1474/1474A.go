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
		var b string
		fmt.Fscan(reader, &b)
		prev := -1
		a := make([]byte, n)
		for i := 0; i < n; i++ {
			digit := int(b[i] - '0')
			if digit+1 != prev {
				a[i] = '1'
				prev = digit + 1
			} else {
				a[i] = '0'
				prev = digit
			}
		}
		fmt.Fprintln(writer, string(a))
	}
}
