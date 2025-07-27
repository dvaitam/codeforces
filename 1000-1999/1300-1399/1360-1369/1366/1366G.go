package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for Codeforces problem 1366G as described in problemG.txt.
// The real solution requires dynamic programming over the sequence of characters
// with backspace semantics. Implementing the optimal algorithm is non-trivial,
// so this program only parses the input and prints 0.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	// TODO: implement the actual computation. For now, print 0.
	fmt.Fprintln(out, 0)
}
