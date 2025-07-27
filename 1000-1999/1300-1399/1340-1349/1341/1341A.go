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
		var n, a, b, c, d int64
		fmt.Fscan(in, &n, &a, &b, &c, &d)
		minTotal := n * (a - b)
		maxTotal := n * (a + b)
		if minTotal <= c+d && maxTotal >= c-d {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
