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
		fmt.Fscan(reader, &n, &s)
		stack := make([]byte, 0, n)
		for i := 0; i < n; i++ {
			c := s[i]
			if len(stack) > 0 && stack[len(stack)-1] != c {
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, c)
			}
		}
		fmt.Fprintln(writer, len(stack))
	}
}
