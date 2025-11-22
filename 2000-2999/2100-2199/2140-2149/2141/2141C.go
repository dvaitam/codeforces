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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// Total operations: 2n (initial push+min for l=0) + n(n-1) for other l
	total := n*n + n
	fmt.Fprintln(out, total)

	// Subarrays starting at 0: after each pushback of a[i], run min on [0..i].
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "pushback a[%d]\n", i)
		fmt.Fprintln(out, "min")
	}

	// For l from 1 to n-1:
	// Drop a[l-1] from the front, then the deque is a[l..n-1].
	// Enumerate subarrays [l..r] in decreasing r by popping back.
	for l := 1; l < n; l++ {
		fmt.Fprintln(out, "popfront")
		for size := n - l; size >= 1; size-- {
			fmt.Fprintln(out, "min")
			if size > 1 {
				fmt.Fprintln(out, "popback")
			}
		}
	}
}
