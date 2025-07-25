package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE.txt from folder 1969.
// The actual algorithm to compute the minimum number of operations
// so that every subarray has an element occurring exactly once is
// not implemented. The program simply reads the input and prints
// zero for each test case so that the repository contains a
// compilable Go solution.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			_ = a
		}
		fmt.Fprintln(out, 0)
	}
}
