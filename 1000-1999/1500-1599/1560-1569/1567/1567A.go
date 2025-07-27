package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to Codeforces problem 1567A - Domino Disaster.
// For each domino placement we are given one row of the 2 x n grid.
// Horizontal dominoes appear the same in both rows (L or R),
// while vertical dominoes swap between U and D.
// Thus the other row is obtained by replacing U with D and vice versa.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		b := []byte(s)
		for i, c := range b {
			if c == 'U' {
				b[i] = 'D'
			} else if c == 'D' {
				b[i] = 'U'
			}
		}
		fmt.Fprintln(out, string(b))
	}
}
