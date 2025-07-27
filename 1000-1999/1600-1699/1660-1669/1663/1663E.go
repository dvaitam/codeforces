package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE.txt in folder 1663.
// The problem statement is incomplete in this repository snapshot.
// This program reads the labyrinth description and produces no output.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(in, &row)
	}
}
