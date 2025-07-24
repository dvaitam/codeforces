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
		var a, b, c, n int64
		fmt.Fscan(in, &a, &b, &c, &n)
		maxv := a
		if b > maxv {
			maxv = b
		}
		if c > maxv {
			maxv = c
		}
		need := (maxv - a) + (maxv - b) + (maxv - c)
		if n < need {
			fmt.Fprintln(out, "NO")
			continue
		}
		if (n-need)%3 == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
