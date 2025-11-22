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
		fmt.Fscan(in, &n)

		// We produce rows that are cyclic shifts of the descending permutation.
		// For the row shifted by k (0-based), reversing the first n-k elements
		// and then the remaining k elements turns the initial row into that shift.
		ops := 2*n - 1
		fmt.Fprintln(out, ops)

		// Row 1: full reverse to get the base descending permutation.
		fmt.Fprintln(out, 1, 1, n)

		for i := 2; i <= n; i++ {
			k := i - 1
			prefix := n - k // length of the first segment
			// Reverse the first prefix elements.
			fmt.Fprintln(out, i, 1, prefix)
			// Reverse the remaining suffix.
			fmt.Fprintln(out, i, prefix+1, n)
		}
	}
}
