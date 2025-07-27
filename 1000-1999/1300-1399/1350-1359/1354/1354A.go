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
		var a, b, c, d int64
		fmt.Fscan(in, &a, &b, &c, &d)
		if b >= a {
			fmt.Fprintln(out, b)
			continue
		}
		if d >= c {
			fmt.Fprintln(out, -1)
			continue
		}
		remain := a - b
		cycle := c - d
		times := (remain + cycle - 1) / cycle
		ans := b + times*c
		fmt.Fprintln(out, ans)
	}
}
