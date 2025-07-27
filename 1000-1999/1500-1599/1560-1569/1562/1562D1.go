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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		var s string
		fmt.Fscan(reader, &s)
		p := make([]int, n+1)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				if s[i] == '+' {
					p[i+1] = p[i] + 1
				} else {
					p[i+1] = p[i] - 1
				}
			} else {
				if s[i] == '+' {
					p[i+1] = p[i] - 1
				} else {
					p[i+1] = p[i] + 1
				}
			}
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			c := p[r] - p[l-1]
			if c == 0 {
				fmt.Fprintln(writer, 0)
			} else if c%2 != 0 {
				fmt.Fprintln(writer, 1)
			} else {
				fmt.Fprintln(writer, 2)
			}
		}
	}
}
