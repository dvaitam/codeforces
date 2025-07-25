package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a minimal placeholder solution for the interactive
// tree reconstruction problem described in problemH.txt.
// The real interactive protocol is not implemented here.
// Instead, the program simply outputs a trivial tree.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 1; i < n; i++ {
		fmt.Fprintf(writer, "%d %d\n", i, i+1)
	}
}
