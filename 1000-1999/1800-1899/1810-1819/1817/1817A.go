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

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		stack := make([]int, 0, r-l)
		for i := l; i < r; i++ {
			stack = append(stack, a[i])
			for len(stack) >= 3 && stack[len(stack)-3] >= stack[len(stack)-2] && stack[len(stack)-2] >= stack[len(stack)-1] {
				last := stack[len(stack)-1]
				stack = append(stack[:len(stack)-2], last)
			}
		}
		fmt.Fprintln(writer, len(stack))
	}
}
