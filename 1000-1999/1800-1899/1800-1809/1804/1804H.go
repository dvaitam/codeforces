package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemH.txt (Codeforces 1804H - Code Lock).
// The actual algorithm for optimizing the wheel arrangement is non-trivial.
// This stub merely parses the input and outputs zeroes.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k, n int
	if _, err := fmt.Fscan(in, &k, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	// Proper implementation should compute the minimal time and the count of
	// optimal arrangements. Here we output zeros as a placeholder.
	fmt.Fprintln(out, 0)
	fmt.Fprintln(out, 0)
}
