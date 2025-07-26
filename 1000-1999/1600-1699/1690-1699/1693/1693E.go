package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE (1693E - unspecified name).
// The real algorithm to minimize the number of operations on array
// was not implemented in this repository snapshot.
//
// This program only reads the input according to the specification
// and outputs 0. It should be replaced with an actual solution if
// needed.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	// TODO: implement the correct algorithm.
	fmt.Fprintln(out, 0)
}
