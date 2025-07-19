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
		var s string
		if _, err := fmt.Fscan(reader, &n, &s); err != nil {
			return
		}
		x := n / 2
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == '(' {
				x += 2
			}
		}
		fmt.Fprintln(writer, x)
	}
}
