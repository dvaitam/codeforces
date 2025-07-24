package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement correct algorithm for rated strings
// Placeholder solution prints 0 for all inputs.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	// No implementation yet
	fmt.Fprintln(writer, 0)
}
