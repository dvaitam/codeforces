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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		stack := make([]int, 0)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == 1 {
				stack = append(stack, 1)
			} else {
				for len(stack) > 0 && stack[len(stack)-1] != x-1 {
					stack = stack[:len(stack)-1]
				}
				if len(stack) > 0 {
					stack[len(stack)-1] = x
				}
			}
			for j, v := range stack {
				if j > 0 {
					writer.WriteByte('.')
				}
				fmt.Fprint(writer, v)
			}
			writer.WriteByte('\n')
		}
	}
}
