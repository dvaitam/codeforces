package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the game described in problemA.txt.
// If the initial number of ones is at least five, Alice can force a win;
// otherwise, Bob wins.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t, n int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(in, &n)
		if n <= 4 {
			fmt.Fprintln(out, "Bob")
		} else {
			fmt.Fprintln(out, "Alice")
		}
	}
}
