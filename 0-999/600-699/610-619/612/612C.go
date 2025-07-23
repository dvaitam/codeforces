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
	fmt.Fscan(reader, &s)

	stack := make([]byte, 0, len(s))
	closing := map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}

	var cnt int
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if _, ok := closing[ch]; ok { // opening bracket
			stack = append(stack, ch)
		} else { // closing bracket
			if len(stack) == 0 {
				fmt.Fprintln(writer, "Impossible")
				return
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if closing[top] != ch {
				cnt++
			}
		}
	}
	if len(stack) != 0 {
		fmt.Fprintln(writer, "Impossible")
		return
	}
	fmt.Fprintln(writer, cnt)
}
