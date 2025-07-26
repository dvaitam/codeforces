package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement solution for problem G from contest 1870.
// This placeholder reads the input format and outputs zeros.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
		}
		// The actual algorithm is not implemented yet.
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, 0)
		}
		fmt.Fprintln(out)
	}
}
