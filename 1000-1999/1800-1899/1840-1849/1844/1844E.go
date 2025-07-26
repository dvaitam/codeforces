package main

import (
	"bufio"
	"fmt"
	"os"
)

// The grid must be colored with letters A, B and C so that
// adjacent cells differ and every 2x2 block contains all
// three letters. Such grids have a simple form: each row
// is a cyclic repetition of the sequence A,B,C shifted by
// some offset, and consecutive rows differ by \u00b11 modulo 3.
//
// A constraint between diagonally adjacent cells (x1,y1)
// and (x2,y2) with x2 = x1+1 determines whether the next
// row's shift is +1 or -1 relative to the current one. We
// only need to verify these \u00b11 relations for every pair of
// adjacent rows. If there is no contradiction, a valid grid
// always exists.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		diff := make([]int8, n) // difference between row i and i+1
		ok := true
		for i := 0; i < k; i++ {
			var x1, y1, x2, y2 int
			fmt.Fscan(in, &x1, &y1, &x2, &y2)
			var sign int8
			if y2 == y1+1 {
				sign = -1
			} else {
				sign = 1
			}
			if diff[x1] == 0 {
				diff[x1] = sign
			} else if diff[x1] != sign {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
