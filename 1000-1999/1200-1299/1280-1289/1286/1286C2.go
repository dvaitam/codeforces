package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a placeholder solution for the interactive problem 1286C2.
// Without an interactive judge available, we simply read n and
// output a fixed string to illustrate the structure of a Go program.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	_ = n
	fmt.Println("a")
}
