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
		var a, b, c int64
		fmt.Fscan(in, &a, &b, &c)

		r := b % 3
		need := (3 - r) % 3
		if c < need {
			fmt.Fprintln(out, -1)
			continue
		}
		c -= need
		tents := a + (b+need)/3 + (c+2)/3
		fmt.Fprintln(out, tents)
	}
}
