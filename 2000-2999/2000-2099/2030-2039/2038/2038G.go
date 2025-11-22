package main

import (
	"bufio"
	"fmt"
	"os"
)

// Offline version of the interactive problem "Guess One Character".
// The input provides the full binary string for each test case, so we can
// simply output any correct position and its character.
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
		var s string
		if _, err := fmt.Fscan(in, &s); err != nil {
			return
		}

		// Choose the first character (1-based index).
		idx := 1
		val := 0
		if len(s) > 0 && s[0] == '1' {
			val = 1
		}
		fmt.Fprintln(out, idx, val)
	}
}
