package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		res := make([]byte, 0, n)
		for i := 0; i < n; {
			c := s[i]
			res = append(res, c)
			i++
			for i < n && s[i] != c {
				i++
			}
			if i < n {
				i++
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
