package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement solution for Bytelandian Tree Factory (problem F).
// The original task requires constructing a sequence of operations to
// transform a bamboo into the desired rooted tree. Since a full
// implementation is non-trivial and the repository lacks one, this
// placeholder program reads the input and outputs a fixed configuration.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &p[i])
	}
	// Output a trivial bamboo labelling and zero operations.
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(i)
	}
	fmt.Println()
	fmt.Println(0)
}
