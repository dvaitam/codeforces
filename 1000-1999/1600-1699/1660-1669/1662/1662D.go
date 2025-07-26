package main

import (
	"bufio"
	"fmt"
	"os"
)

func canonical(s string) (int, string) {
	b := 0
	stack := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == 'B' {
			b ^= 1
		} else {
			if len(stack) > 0 && stack[len(stack)-1] == ch {
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, ch)
			}
		}
	}
	return b, string(stack)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var u, v string
		fmt.Fscan(reader, &u)
		fmt.Fscan(reader, &v)
		bu, su := canonical(u)
		bv, sv := canonical(v)
		if bu == bv && su == sv {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
