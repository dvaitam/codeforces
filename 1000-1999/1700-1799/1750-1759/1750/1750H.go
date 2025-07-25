package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the algorithm for counting good substrings as described in problemH.txt.
// This placeholder only parses the input and outputs zero for each test case.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		// Placeholder implementation
		fmt.Fprintln(writer, 0)
	}
}
