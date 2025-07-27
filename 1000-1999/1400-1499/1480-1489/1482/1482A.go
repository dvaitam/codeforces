package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// The prison has a grid of a*b cells. To ensure an escape path
// from every cell with the minimum number of broken walls, we
// can connect all cells into one tree and open one wall to the
// outside. Any such spanning tree contains exactly a*b edges,
// so breaking a*b walls is both necessary and sufficient.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b int
		fmt.Fscan(in, &a, &b)
		fmt.Fprintln(out, a*b)
	}
}
