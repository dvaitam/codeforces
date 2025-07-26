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
		var n, m int
		fmt.Fscan(in, &n, &m)
		switch {
		case n == 1 && m == 1:
			fmt.Fprintln(out, 0)
		case n == 1 || m == 1:
			fmt.Fprintln(out, 1)
		default:
			fmt.Fprintln(out, 2)
		}
	}
}
