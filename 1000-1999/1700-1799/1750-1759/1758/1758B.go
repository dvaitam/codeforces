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
		switch {
		case n == 1:
			fmt.Fprintln(writer, 69)
		case n == 2:
			fmt.Fprintln(writer, "1 3")
		case n%2 == 1:
			for i := 0; i < n; i++ {
				writer.WriteString("7")
				if i < n-1 {
					writer.WriteByte(' ')
				}
			}
			writer.WriteByte('\n')
		default:
			for i := 0; i < n-3; i++ {
				writer.WriteString("2 ")
			}
			writer.WriteString("1 2 3\n")
		}
	}
}
