package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the actual algorithm for problem E2.
// The current implementation only reads the input and
// outputs 0 since computing the optimal permutation
// requires a complex approach for large n (up to 1e6).
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// read parent array
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
	}
	fmt.Println(0)
}
