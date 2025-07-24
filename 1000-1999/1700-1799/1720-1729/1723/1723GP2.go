package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	// read edges, ignore contents
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
	}
	var setSize int
	var k float64
	fmt.Fscan(in, &setSize, &k)
	for i := 0; i < setSize; i++ {
		var x int
		fmt.Fscan(in, &x)
	}

	// Output a trivial partition: single set containing all nodes
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, 1)
	fmt.Fprintf(out, "%d", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, " %d", i)
	}
	fmt.Fprintln(out)
}
