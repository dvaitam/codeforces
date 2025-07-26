package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution.
// This placeholder reads the input and outputs zeros for each test case.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, 0)
	}
}
