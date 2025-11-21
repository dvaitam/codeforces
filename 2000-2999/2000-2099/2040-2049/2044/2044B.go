package main

import (
	"bufio"
	"fmt"
	"os"
)

func transform(s string) string {
	n := len(s)
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		ch := s[n-1-i]
		switch ch {
		case 'p':
			res[i] = 'q'
		case 'q':
			res[i] = 'p'
		default: // 'w'
			res[i] = ch
		}
	}
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a string
		fmt.Fscan(in, &a)
		fmt.Fprintln(out, transform(a))
	}
}
