package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// For each test case we are given integers l and r and need
// to maximize a mod b where r >= a >= b >= l. The best choice
// for a is a = r. If we can pick b <= r/2 + 1, then r mod b is
// maximized when b = r/2 + 1. Otherwise b must be at least l,
// and since b > r/2, r mod b simplifies to r - b, maximized at
// b = l. Therefore the answer is r mod max(l, r/2+1).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		b := r/2 + 1
		if l > b {
			b = l
		}
		fmt.Fprintln(out, r%b)
	}
}
