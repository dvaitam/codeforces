package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(s string) string {
	b := []byte(s)
	n := len(b)
	// propagate known characters to the right
	for i := 1; i < n; i++ {
		if b[i] == '?' && b[i-1] != '?' {
			b[i] = b[i-1]
		}
	}
	// fill remaining '?' from right to left
	for i := n - 1; i >= 0; i-- {
		if b[i] == '?' {
			if i+1 < n {
				b[i] = b[i+1]
			} else {
				b[i] = '0'
			}
		}
	}
	return string(b)
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
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, solve(s))
	}
}
