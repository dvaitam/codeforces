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
		var s string
		fmt.Fscan(reader, &s)
		last := make([]int, 26)
		for i := range last {
			last[i] = -1
		}
		for i := 0; i < len(s); i++ {
			last[int(s[i]-'a')] = i
		}
		used := make([]bool, 26)
		stack := make([]byte, 0, 26)
		for i := 0; i < len(s); i++ {
			c := s[i]
			idx := int(c - 'a')
			if used[idx] {
				continue
			}
			for len(stack) > 0 && stack[len(stack)-1] < c && last[int(stack[len(stack)-1]-'a')] > i {
				used[int(stack[len(stack)-1]-'a')] = false
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)
			used[idx] = true
		}
		fmt.Fprintln(writer, string(stack))
	}
}
