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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	stack := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		stack = append(stack, x)
		for len(stack) >= 2 && stack[len(stack)-1] == stack[len(stack)-2] {
			val := stack[len(stack)-1] + 1
			stack = stack[:len(stack)-2]
			stack = append(stack, val)
		}
	}
	fmt.Fprintln(writer, len(stack))
	for i, v := range stack {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	if len(stack) > 0 {
		fmt.Fprintln(writer)
	}
}
