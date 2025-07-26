package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program checks if each input string is square, meaning it is
// some string written twice in a row.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		if len(s)%2 == 0 && s[:len(s)/2] == s[len(s)/2:] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
