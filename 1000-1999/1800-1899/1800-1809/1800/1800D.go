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
		dup := 0
		for i := 0; i < n-2; i++ {
			if s[i] == s[i+2] {
				dup++
			}
		}
		fmt.Fprintln(writer, n-1-dup)
	}
}
