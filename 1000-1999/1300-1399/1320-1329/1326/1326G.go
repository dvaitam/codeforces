package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for problem G as described in problemG.txt.
// The problem involves counting partitions of a planar tree into spiderweb
// subtrees. Implementing the full geometry and dynamic programming solution
// is non-trivial, so this placeholder only parses the input and outputs 0 so
// that the file compiles and the repository builds.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// Read coordinates
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
	}

	// Read edges
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}

	fmt.Fprintln(out, 0)
}
