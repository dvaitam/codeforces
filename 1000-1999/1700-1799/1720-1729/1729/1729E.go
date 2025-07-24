package main

import (
	"bufio"
	"fmt"
	"os"
)

// The interactive version asks to determine the size of a hidden cycle.
// In the non-interactive archive version used for hacking, the input
// simply contains the value of n.  We read this value and output it
// directly.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fprintln(out, n)
}
