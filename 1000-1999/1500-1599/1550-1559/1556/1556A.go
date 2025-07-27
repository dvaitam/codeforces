package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt
// for contest 1556A (A Variety of Operations). We start from (0,0)
// and can perform operations that either add the same k to both
// numbers or add/subtract k from them in opposite directions.
//
// Observations:
//  1. The parity of c and d must match. All operations preserve the
//     parity of a+b and a-b, which are initially even. Thus, if c+d
//     is odd, the target is unreachable.
//  2. If c and d are both zero, no operations are needed.
//  3. If c equals d (and is non-zero), a single operation of the first
//     type with k=c suffices.
//  4. In all other reachable cases two operations are enough:
//     add (c+d)/2 to both numbers and then adjust their difference
//     using an operation of the second or third type.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var c, d int
		fmt.Fscan(in, &c, &d)
		if c == 0 && d == 0 {
			fmt.Fprintln(out, 0)
		} else if (c+d)%2 == 1 {
			fmt.Fprintln(out, -1)
		} else if c == d {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
