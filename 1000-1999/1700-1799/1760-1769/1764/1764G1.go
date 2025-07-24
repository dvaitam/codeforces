package main

import (
	"bufio"
	"fmt"
	"os"
)

// Note: This file contains only a minimal interactive template.
// The actual strategy for locating the index y such that p_y = 1 is
// not implemented as the full problem statement was unavailable.
// The program simply outputs 1 without issuing queries.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// The intended solution should perform interactive queries of the form:
	// fmt.Fprintf(out, "? %d %d %d\n", l, r, k)
	// and then read the integer response via fmt.Fscan(in, &ans).
	// Using at most 30 such queries, the algorithm should deduce the
	// position of the value 1 in the hidden permutation.

	// Placeholder behaviour: output 1 directly.
	fmt.Fprintln(out, "! 1")
}
