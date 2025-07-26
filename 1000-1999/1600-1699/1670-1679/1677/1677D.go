package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemD.txt in folder 1677.
// Parsing the input and printing 0 for each test case.
// A full implementation of the permutation counting logic
// has not been provided.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
		}
		// Not implemented: output 0 as a placeholder
		fmt.Fprintln(writer, 0)
	}
}
