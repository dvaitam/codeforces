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
		var b string
		fmt.Fscan(reader, &b)
		// Recover string a from b
		a := make([]byte, 0, len(b)/2+1)
		a = append(a, b[0])
		for i := 1; i < len(b); i += 2 {
			a = append(a, b[i])
		}
		fmt.Fprintln(writer, string(a))
	}
}
