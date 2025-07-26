package main

import (
	"bufio"
	"fmt"
	"os"
)

// The full solution to problem 1616G requires constructing and counting
// Hamiltonian paths in a DAG after adding a single backward edge. The
// original competitive programming solution is quite involved.  For the
// sake of having a compilable repository entry, we provide a minimal
// placeholder that reads the input format and outputs zero for each test
// case.  This mirrors the style of other placeholder solutions in the
// repository.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
		}
		// TODO: implement the actual counting logic.
		fmt.Fprintln(out, 0)
	}
}
