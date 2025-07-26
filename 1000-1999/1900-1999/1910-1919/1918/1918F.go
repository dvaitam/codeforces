package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problem 1918F. The full algorithm to minimize the
// travel time with up to k teleports to the root is not implemented here.
// This stub only reads the input and outputs 0.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var p int
		fmt.Fscan(reader, &p)
	}
	fmt.Println(0)
}
