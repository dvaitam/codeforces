package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverseString(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
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
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		rev := reverseString(s)
		ans := s
		if tmp := s + rev; tmp < ans {
			ans = tmp
		}
		if tmp := rev + s; tmp < ans {
			ans = tmp
		}
		fmt.Fprintln(out, ans)
	}
}
