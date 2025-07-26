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
		var s string
		fmt.Fscan(reader, &s)
		ops := 0
		for i := 0; i < n; i += 2 {
			if s[i] != s[i+1] {
				ops++
			}
		}
		fmt.Fprintln(writer, ops)
	}
}
