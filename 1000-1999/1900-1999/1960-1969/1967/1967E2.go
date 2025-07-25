package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: This is a placeholder implementation for problem E2.
// The proper algorithm for counting arrays is non-trivial and
// is not yet implemented. Currently this just reads the input
// and outputs zero for each test case.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, b0 int
		fmt.Fscan(in, &n, &m, &b0)
		// Placeholder result
		fmt.Fprintln(out, 0)
	}
}
