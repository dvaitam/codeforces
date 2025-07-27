package main

import (
	"bufio"
	"fmt"
	"os"
)

// Interactive solution for Codeforces problem 1583D - Omkar and the Meaning of Life.
// The program issues 2n queries to discover the hidden permutation. It first
// queries arrays with a single element increased to find successors with
// a larger index. Then it queries arrays with a single element decreased to
// find successors with a smaller index and also locate the position of 1.
// Finally the permutation is reconstructed and printed.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	succ := make([]int, n+1)
	pos1 := 0

	// First round of n queries: use [1,1,...,2,...,1]
	for i := 1; i <= n; i++ {
		fmt.Fprint(out, "?")
		for j := 1; j <= n; j++ {
			val := 1
			if j == i {
				val = 2
			}
			fmt.Fprintf(out, " %d", val)
		}
		fmt.Fprintln(out)
		out.Flush()

		var k int
		if _, err := fmt.Fscan(in, &k); err != nil {
			return
		}
		if k != 0 && k < i {
			succ[k] = i
		}
	}

	// Second round of n queries: use [2,2,...,1,...,2]
	for i := 1; i <= n; i++ {
		fmt.Fprint(out, "?")
		for j := 1; j <= n; j++ {
			val := 2
			if j == i {
				val = 1
			}
			fmt.Fprintf(out, " %d", val)
		}
		fmt.Fprintln(out)
		out.Flush()

		var k int
		if _, err := fmt.Fscan(in, &k); err != nil {
			return
		}
		if k == 0 {
			pos1 = i
		} else if k < i {
			succ[k] = i
		}
	}

	// Reconstruct permutation starting from the index of value 1
	perm := make([]int, n+1)
	idx := pos1
	for v := 1; v <= n; v++ {
		perm[idx] = v
		idx = succ[idx]
	}

	fmt.Fprint(out, "!")
	for i := 1; i <= n; i++ {
		fmt.Fprintf(out, " %d", perm[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
