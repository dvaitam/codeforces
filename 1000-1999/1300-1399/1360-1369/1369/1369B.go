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
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		prefix := 0
		for prefix < n && s[prefix] == '0' {
			prefix++
		}
		suffix := n - 1
		for suffix >= 0 && s[suffix] == '1' {
			suffix--
		}
		if prefix > suffix {
			fmt.Fprintln(writer, s)
		} else {
			result := s[:prefix] + "0" + s[suffix+1:]
			fmt.Fprintln(writer, result)
		}
	}
}
