package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a placeholder solution for the interactive problem described in
// problemC.txt. The real challenge requires playing a game where Alice and Bob
// alternately add and remove numbers from a set in order to maximize or
// minimize its final MEX. Since this repository does not provide an interactive
// judge, the program merely reads the provided input format and exits without
// performing any interactive steps.
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
		// Read the initial set S of size n.
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
		}
		// Normally Alice would start interacting with the judge here.
		// Without an interactive environment we stop after parsing input.
	}
}
