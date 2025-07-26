package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: This is a placeholder solution for problem D.
// The actual algorithm for finding the minimal number of tricks
// is non-trivial and has not been implemented yet. Currently the
// program just reads all input and outputs 0 for each test case.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}
		// Placeholder result
		fmt.Fprintln(out, 0)
	}
}
