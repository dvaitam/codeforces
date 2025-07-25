package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a placeholder implementation for problem H1 from contest 1781.
// The actual counting logic for distinct light configurations has not
// been implemented yet. The program reads the input in the required
// format and outputs zero for each test case so that it compiles and
// can be expanded later.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var h, w, k int
		fmt.Fscan(in, &h, &w, &k)
		for i := 0; i < k; i++ {
			var r, c int
			fmt.Fscan(in, &r, &c)
		}
		// TODO: implement the actual logic for the easy version.
		fmt.Fprintln(out, 0)
	}
}
