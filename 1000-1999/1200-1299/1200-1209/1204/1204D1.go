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

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	n := len(s)
	res := []byte(s)
	stack := make([]int, 0)
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			stack = append(stack, i)
		} else {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	for _, idx := range stack {
		res[idx] = '0'
	}
	fmt.Fprintln(writer, string(res))
}
