package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the actual algorithm for problem E2. The constraints
// allow n up to 2e5 and a_i up to 1e9. The current implementation is a
// placeholder that simply reads the input and outputs zero monsters for
// every test case.
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
			var x int64
			fmt.Fscan(in, &x)
		}
		fmt.Fprintln(out, 0)
	}
}
