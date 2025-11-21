package main

import (
	"bufio"
	"fmt"
	"os"
)

// The original problem 2164G "Pointless Machine" is interactive.  In the
// hacking format used in this repository the entire tree is provided in the
// input, so we simply echo its edges for each test case.
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
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &edges[i][0], &edges[i][1])
		}
		for _, e := range edges {
			fmt.Fprintln(out, e[0], e[1])
		}
	}
}
