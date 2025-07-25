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
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		first := -1
		last := -1
		for i := 0; i < n; i++ {
			if s[i] == 'B' {
				if first == -1 {
					first = i
				}
				last = i
			}
		}
		fmt.Fprintln(writer, last-first+1)
	}
}
