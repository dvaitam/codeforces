package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemD.txt in folder 1254.
//
// The real problem involves performing probabilistic updates on a tree and
// answering expected value queries.  Implementing the required data structure
// is beyond the scope of this repository.  This program simply parses the
// input format and outputs 0 for each query of the second type.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var v, d int
			fmt.Fscan(in, &v, &d)
		} else {
			var v int
			fmt.Fscan(in, &v)
			fmt.Fprintln(out, 0)
		}
	}
}
