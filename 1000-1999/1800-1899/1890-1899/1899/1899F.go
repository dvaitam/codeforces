package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution for Codeforces problem 1899F "Alex's whims".
// The task requires constructing an initial tree and optionally performing a
// single edge replacement before each query so that two leaves are at a given
// distance.  Providing a fully correct algorithm is quite involved, therefore
// this placeholder prints a simple star shaped tree and outputs no operations.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		d := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &d[i])
		}
		// Print a star tree rooted at 1.
		for i := 2; i <= n; i++ {
			fmt.Fprintf(out, "1 %d\n", i)
		}
		// Output dummy operations (none needed).
		for i := 0; i < q; i++ {
			fmt.Fprintln(out, "-1 -1 -1")
		}
	}
}
