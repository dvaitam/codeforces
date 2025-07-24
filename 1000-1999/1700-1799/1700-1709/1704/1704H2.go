package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a placeholder implementation for problem H2 from contest 1704.
// The real solution requires a combinatorial analysis which is currently
// not implemented. The program merely reads the input values and prints
// zeroes so it compiles and can serve as a starting point for future
// work.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k, M int
	if _, err := fmt.Fscan(in, &k, &M); err != nil {
		return
	}

	// TODO: implement the actual counting logic.
	for i := 1; i <= k; i++ {
		fmt.Fprintln(out, 0)
	}
}
