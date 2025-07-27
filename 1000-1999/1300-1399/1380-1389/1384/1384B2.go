package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Placeholder solution. The full implementation of the hard version
	// of problem B2 is not provided.
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k, l int
		fmt.Fscan(in, &n, &k, &l)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
		}
		fmt.Fprintln(out, "No")
	}
}
