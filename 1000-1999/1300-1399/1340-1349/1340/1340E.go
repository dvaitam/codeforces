package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for interactive problem 1340E ("Three strange Bees").
// The real interactive strategy is not implemented. This program simply reads
// the graph description so that it compiles and outputs a constant value.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	// TODO: implement strategy to place and move bees to catch Nastya.
	fmt.Fprintln(out, 0)
}
