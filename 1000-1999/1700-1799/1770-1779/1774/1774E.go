package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, d int
	if _, err := fmt.Fscan(reader, &n, &d); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
	}
	var m1 int
	fmt.Fscan(reader, &m1)
	for i := 0; i < m1; i++ {
		var x int
		fmt.Fscan(reader, &x)
	}
	var m2 int
	fmt.Fscan(reader, &m2)
	for i := 0; i < m2; i++ {
		var x int
		fmt.Fscan(reader, &x)
	}

	// TODO: implement algorithm to compute minimum steps while keeping the
	// distance between the two pieces not exceeding d. The pieces start at
	// node 1, must visit their respective sets of nodes in any order, and
	// finally return to the root.
	// Placeholder implementation outputs 0.
	fmt.Fprintln(writer, 0)
}
