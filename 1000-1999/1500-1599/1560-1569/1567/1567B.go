package main

import (
	"bufio"
	"fmt"
	"os"
)

func prefixXor(n int) int {
	switch n & 3 {
	case 0:
		return n
	case 1:
		return 1
	case 2:
		return n + 1
	default:
		return 0
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b int
		fmt.Fscan(reader, &a, &b)

		x := prefixXor(a - 1)
		if x == b {
			fmt.Fprintln(writer, a)
		} else if (x ^ b) == a {
			fmt.Fprintln(writer, a+2)
		} else {
			fmt.Fprintln(writer, a+1)
		}
	}
}
