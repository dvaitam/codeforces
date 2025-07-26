package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA from contest 1632.
// We are given a binary string s. We may reorder its
// characters arbitrarily and want a permutation that
// contains no palindromic substring of length greater
// than one.
// It is known that any binary string of length at least
// three will have a palindromic substring of length two
// or three. For length two, it happens when two adjacent
// characters are equal. If all adjacent characters are
// different, the pattern becomes alternating and forms a
// palindrome of length three ("010" or "101"). Therefore
// such strings exist only for lengths one and two with
// distinct characters.
// Hence, the answer is YES if n == 1 or (n == 2 and the
// two characters are different); otherwise, it is NO.
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
		var s string
		fmt.Fscan(in, &n, &s)
		if n == 1 {
			fmt.Fprintln(out, "YES")
		} else if n == 2 && s[0] != s[1] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
