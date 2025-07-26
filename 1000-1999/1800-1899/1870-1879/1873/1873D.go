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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		ops := 0
		i := 0
		for i < n {
			if s[i] == 'B' {
				ops++
				i += k
			} else {
				i++
			}
		}
		fmt.Fprintln(writer, ops)
	}
}
