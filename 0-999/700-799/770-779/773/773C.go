package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder implementation for problem 773C (Prairie Partition).
// The full algorithm is non-trivial and is not implemented here.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
	}
	fmt.Println(-1)
}
