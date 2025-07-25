package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the solution for problem B.
// The current implementation only reads the input and outputs "NO" for
// each test case. The actual algorithm for this problem is non-trivial and
// has not been implemented yet.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var N int
		var A, B string
		fmt.Fscan(in, &N)
		fmt.Fscan(in, &A)
		fmt.Fscan(in, &B)
		// TODO: solve for each test case
		fmt.Fprintln(out, "NO")
	}
}
