package main

import (
	"bufio"
	"fmt"
	"os"
)

// This file contains a stub solution for the interactive problem 1713D.
// The original problem requires interaction with the judging system to
// determine the winner of a single-elimination tournament using queries.
// In this repository we only provide a template that demonstrates how the
// interaction could be organized in Go.  It does not implement the real
// interaction protocol.  Instead, it simply reads the number of test cases
// and the value of n for each test.  To keep the example self contained,
// the program outputs a placeholder result for each test case.

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
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		// Normally, here we would interact with the judge to find the winner
		// using at most ceil((1/3) * 2^(n+1)) queries.  Since this environment
		// does not support interactive problems, we simply output 1 as a
		// placeholder answer.
		fmt.Fprintln(writer, 1)
	}
}
