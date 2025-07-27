package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problem B from Codeforces contest 1526.
// For each query x, it checks whether x can be expressed as a sum of numbers
// consisting solely of digit 1 and having length at least two (11, 111, ...).
// The key observation is that x is representable iff 111*(x%11) <= x.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x int
		fmt.Fscan(in, &x)
		if 111*(x%11) <= x {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
