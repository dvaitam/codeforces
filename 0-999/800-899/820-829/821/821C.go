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

	stack := make([]int, 0)
	expected := 1
	reorders := 0

	for i := 0; i < 2*n; i++ {
		var cmd string
		fmt.Fscan(reader, &cmd)
		if cmd == "add" {
			var x int
			fmt.Fscan(reader, &x)
			stack = append(stack, x)
		} else { // remove
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == expected {
					stack = stack[:len(stack)-1]
				} else {
					reorders++
					stack = stack[:0]
				}
			}
			expected++
		}
	}

	fmt.Fprintln(writer, reorders)
}
