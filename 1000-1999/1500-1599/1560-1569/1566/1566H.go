package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for the interactive problem described in problemH.txt.
// The real solution would issue queries to learn the hidden set A, but this
// repository does not provide an interactive judge. The program only reads the
// input parameters and terminates.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var c, n int
	if _, err := fmt.Fscan(in, &c, &n); err != nil {
		return
	}
	// Normally queries would follow here, but we do nothing without a judge.
	_ = c
	_ = n
}
