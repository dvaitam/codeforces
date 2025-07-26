package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt (GCD Arrays).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var l, r, k int64
		fmt.Fscan(in, &l, &r, &k)

		if l == r {
			if l > 1 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		// Count odd numbers in [l, r].
		odd := (r+1)/2 - l/2
		if k >= odd {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
