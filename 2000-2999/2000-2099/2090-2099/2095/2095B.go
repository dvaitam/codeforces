package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// The answer does not depend on the input; any column outside [0,16] wins.
	// See editorial for 2095B (April Fools Day 2025).
	in := bufio.NewReader(os.Stdin)
	var dummy string
	// Consume the only line if present to match the required I/O format.
	in.ReadString('\n')
	_ = dummy
	fmt.Println("-10000")
}
