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
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		mask := n - 1
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			if (i & (mask - i)) == 0 {
				fmt.Fprint(out, k)
			} else {
				fmt.Fprint(out, 0)
			}
		}
		fmt.Fprintln(out)
	}
}
