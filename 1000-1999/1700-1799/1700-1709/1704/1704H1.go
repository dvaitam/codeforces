package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for problem H1.
// The current implementation only reads k and M and prints 0.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var k, M int
	if _, err := fmt.Fscan(reader, &k, &M); err != nil {
		return
	}

	// Placeholder implementation.
	fmt.Fprintln(writer, 0)
}
