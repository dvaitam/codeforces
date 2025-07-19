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
		i := 1
		for i < len(s) && s[i-1] != s[i] {
			i++
		}
		ch := byte('a')
		if s[i-1] == 'a' {
			ch = 'b'
		}
		if i == len(s) {
			s += string(ch)
		} else {
			s = s[:i] + string(ch) + s[i:]
		}
		fmt.Fprintln(writer, s)
	}
}
