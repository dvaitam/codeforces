package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the exact algorithm that computes the probability described
// in problemH.txt. This placeholder only parses the input and outputs 0.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	for i := 0; i < m; i++ {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		_ = t
		_ = x
		_ = y
	}

	// Placeholder output
	fmt.Fprintln(writer, 0)
}
