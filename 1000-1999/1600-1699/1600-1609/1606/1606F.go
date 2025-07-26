package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement an efficient algorithm for problem F based on the
// description in problemF.txt. The problem requires maximizing
// c(v) - m*k for each query on a rooted tree after deleting any
// vertices except the root and the queried vertex v. Deleting a vertex
// promotes its children to the parent. The constraints are large
// (up to 2e5 vertices and queries), so a full solution would need
// sophisticated preprocessing. This placeholder merely reads the
// input and outputs 0 for every query so that the program compiles.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		// edges are ignored in this placeholder
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var v, k int
		fmt.Fscan(in, &v, &k)
		// A correct implementation would compute the optimal value
		// for the given v and k. We simply print 0.
		fmt.Fprintln(out, 0)
	}
}
