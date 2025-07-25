package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE.txt (Happy Life in University).
// The actual algorithm is not implemented. This program only
// reads the input according to the specification and outputs
// zero for each test case so that the build succeeds.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 2; i <= n; i++ {
			var p int
			fmt.Fscan(in, &p)
			_ = p
		}
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(in, &a)
			_ = a
		}
		fmt.Fprintln(out, 0)
	}
}
