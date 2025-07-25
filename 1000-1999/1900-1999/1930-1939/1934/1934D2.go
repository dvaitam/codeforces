package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// This file contains a stub solution for the interactive problem D2 of
// contest 1934.  The actual task requires interacting with a judge, but
// this repository does not provide such a judge.  The program therefore
// demonstrates the basic structure of a solution and prints a single
// move based on the parity of the number of set bits in n.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n uint64
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		if bits.OnesCount64(n)%2 == 0 {
			// Even popcount is a winning state, choose to move first.
			fmt.Fprintln(out, "first")
			out.Flush()

			// Use the least significant set bit for the initial split.
			p1 := n & -n
			p2 := n - p1
			if p2 == 0 {
				p1, p2 = 1, n-1
			}
			fmt.Fprintf(out, "%d %d\n", p1, p2)
			out.Flush()
			// A real interactive solution would continue reading
			// Bob's choice and responding until the game ends.
			return
		} else {
			// Odd popcount is losing, so go second.
			fmt.Fprintln(out, "second")
			out.Flush()
			return
		}
	}
}
