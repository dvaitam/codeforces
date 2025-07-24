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

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	// count plus and minus signs
	p := 0
	for i := 0; i < n; i++ {
		if s[i] == '+' {
			p++
		}
	}
	m := n - p

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if a == b {
			if p == m {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}
		numer := b * int64(m-p)
		denom := a - b
		// check divisibility
		if numer%denom != 0 {
			fmt.Fprintln(out, "NO")
			continue
		}
		sa := numer / denom
		if sa < int64(-m) || sa > int64(p) {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
