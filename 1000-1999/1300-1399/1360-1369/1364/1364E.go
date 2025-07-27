package main

import (
	"bufio"
	"fmt"
	"os"
)

// This file contains a stub solution for the interactive problem 1364E.
// The real problem asks to determine a hidden permutation of length n by
// issuing queries that return the bitwise OR of two elements. In this
// repository there is no interactive judge available, so the program
// simply reads n and prints the identity permutation as a placeholder.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, i)
	}
	writer.WriteByte('\n')
}
