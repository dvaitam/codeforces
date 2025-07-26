package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the expected number of moves as described in problemH.txt.
// This placeholder only parses the input and outputs 0 for each test case.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		// Proper implementation should compute the expected moves modulo 998244353.
		fmt.Fprintln(out, 0)
	}
}
