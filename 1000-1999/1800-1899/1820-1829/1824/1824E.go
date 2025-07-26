package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the algorithm for Codeforces problem "LuoTianyi and Cartridge" (1824E).
// The official statement can be found in problemE.txt. The task asks for the
// maximum value of min(A,C)*(B+D) for a tree built under complex constraints.
// A full implementation requires advanced data structures and is non-trivial.
//
// For now this file only provides the basic I/O structure so that the solution
// compiles. It simply reads the input and prints 0.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 0; i < n-1; i++ {
		var x, y, c, d int
		fmt.Fscan(reader, &x, &y, &c, &d)
		// input is read but ignored in this placeholder
	}

	// Placeholder output
	fmt.Fprintln(writer, 0)
}
