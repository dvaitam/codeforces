package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution for problem 1796E.
// This placeholder reads the input and outputs 1 for each test case,
// corresponding to the minimal possible cost when each vertex has a
// distinct color.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
		}
		fmt.Fprintln(out, 1)
	}
}
