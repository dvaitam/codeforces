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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, a, q int
		fmt.Fscan(in, &n, &a, &q)
		var s string
		fmt.Fscan(in, &s)
		cur := a
		maxOnline := a
		plus := 0
		for _, ch := range s {
			if ch == '+' {
				cur++
				plus++
			} else {
				cur--
			}
			if cur > maxOnline {
				maxOnline = cur
			}
		}
		if maxOnline >= n {
			fmt.Fprintln(out, "YES")
		} else if a+plus < n {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "MAYBE")
		}
	}
}
