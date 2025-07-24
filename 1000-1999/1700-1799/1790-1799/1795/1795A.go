package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It determines whether two towers of blocks can be rearranged
// using the allowed operations so that within each tower no
// adjacent blocks share the same color.

func countDup(s string) int {
	c := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			c++
		}
	}
	return c
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s, u string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &u)

		ds := countDup(s)
		du := countDup(u)

		if ds+du > 1 || (ds+du == 1 && s[len(s)-1] == u[len(u)-1]) {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
