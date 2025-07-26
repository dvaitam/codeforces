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
		var b, c, h int
		fmt.Fscan(in, &b, &c, &h)
		f := c + h
		x := b
		if f+1 < b {
			x = f + 1
		}
		fmt.Fprintln(out, 2*x-1)
	}
}
