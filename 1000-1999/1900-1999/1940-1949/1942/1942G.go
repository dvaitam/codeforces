package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the real solution for problem G as described in problemG.txt.
// Currently this is just a placeholder that parses the input format and prints 0
// for each test case so that the repository builds and runs.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		// Placeholder output
		fmt.Fprintln(writer, 0)
	}
}
