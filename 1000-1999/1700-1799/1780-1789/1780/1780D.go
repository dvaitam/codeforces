package main

import (
	"bufio"
	"fmt"
	"os"
)

// The original Codeforces problem 1780D is an interactive problem.
// The statement provides the initial popcount of a hidden number and
// allows up to 30 queries where a positive integer x is subtracted
// from the hidden number and the new popcount is returned.
//
// In this repository we do not have an interactive judge, so the
// implementation below is only a placeholder illustrating how the
// interaction could be organised.  Without the interactor the program
// simply reads the number of test cases and the initial popcount for
// each case and terminates without further interaction.
//
// When solving on Codeforces one would print queries in the form
// "? x" (or "- x" depending on the exact protocol), flush the output
// and read the updated popcount.  After the value of n is determined
// it should be printed with the prefix "!".
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var cnt int
		if _, err := fmt.Fscan(in, &cnt); err != nil {
			return
		}
		// TODO: interactive protocol is not implemented.
		// This placeholder does nothing further because the
		// repository does not provide an interactor.
		_ = cnt
	}
}
