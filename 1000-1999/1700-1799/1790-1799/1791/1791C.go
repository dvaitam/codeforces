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
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		l, r := 0, n-1
		for l < r && s[l] != s[r] {
			l++
			r--
		}
		if r < l {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, r-l+1)
		}
	}
}
