package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Provide the actual game theoretic solution for problem G.
// This placeholder implementation merely counts the initial number of
// black cells and decides the winner by its parity.  The real problem
// requires analysing how the toggle operation interacts between cells.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var r, c int
		fmt.Fscan(reader, &r, &c)
		_ = r
		_ = c
	}

	if n%2 == 1 {
		fmt.Fprintln(writer, "FIRST")
	} else {
		fmt.Fprintln(writer, "SECOND")
	}
}
