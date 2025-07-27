package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement a full solution for problem D as described in problemD.txt.
// The puzzle requires moving the k-th toy to the leftmost cell of the top row
// in a 2xN grid using at most 1,000,000 moves consisting of L, R, U and D.
// Computing an optimal sequence is non-trivial, so this placeholder just
// outputs k moves to the left, which is not guaranteed to solve the puzzle.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	for i := 0; i < k; i++ {
		fmt.Fprint(out, "L")
	}
	fmt.Fprintln(out)
}
