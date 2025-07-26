package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for interactive problem 1867E1.
// The actual interactive strategy is not implemented.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		// TODO: implement interactive strategy to compute XOR of the array
		fmt.Fprintln(out, 0)
	}
}
