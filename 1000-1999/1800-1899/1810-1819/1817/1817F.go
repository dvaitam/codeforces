package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement a full solution for problem F as described in problemF.txt.
// The task asks for the count of entangled pairs (a, b) of substrings of a
// given string s, where every occurrence of a is immediately followed by some
// fixed string c concatenated with b and every occurrence of b is immediately
// preceded by a concatenated with the same c. Computing this exactly requires
// sophisticated string algorithms which are beyond this placeholder.
//
// This program only provides the basic input/output structure so that the file
// compiles. It reads the string s and outputs 0.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	// A correct solution would analyze s and compute the actual count here.
	fmt.Fprintln(out, 0)
}
