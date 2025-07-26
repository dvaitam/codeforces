package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the correct algorithm for problem E as described in problemE.txt.
// The current implementation only parses the input and outputs zero for each test case.
// It serves as a placeholder.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k int
		var s string
		fmt.Fscan(in, &k, &s)
		// Placeholder output
		fmt.Fprintln(out, 0)
	}
}
