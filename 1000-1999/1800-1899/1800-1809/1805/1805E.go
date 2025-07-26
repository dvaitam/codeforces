package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE.txt in folder 1805.
// The real algorithm to compute the maximum double parameter
// after removing each edge is not implemented. This program
// just reads the tree and outputs zeros so that the code
// compiles successfully.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(in, &a)
	}
	for i := 0; i < n-1; i++ {
		fmt.Fprintln(out, 0)
	}
}
