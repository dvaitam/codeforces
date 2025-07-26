package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution: always output 1 segment per test case.
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
		fmt.Fprintln(out, 1)
	}
}
