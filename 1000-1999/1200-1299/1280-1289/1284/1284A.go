package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It generates the name of a year by cycling through two
// sequences of strings and concatenating the corresponding
// elements.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	t := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &t[i])
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var y int
		fmt.Fscan(in, &y)
		fmt.Fprintln(out, s[(y-1)%n]+t[(y-1)%m])
	}
}
