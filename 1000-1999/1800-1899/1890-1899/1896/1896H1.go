package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for problem H1.
// The current implementation only reads input and prints 0.
// This is a placeholder so the file compiles.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	// Placeholder implementation.
	fmt.Fprintln(out, 0)
}
