package main

import (
	"bufio"
	"fmt"
	"os"
)

var rdr = bufio.NewReader(os.Stdin)

func readInt64() int64 {
	var x int64
	fmt.Fscan(rdr, &x)
	return x
}

func main() {
	// Placeholder solution for problem G.
	// Read a single integer (ignored) and output 0.
	_ = readInt64()
	fmt.Println(0)
}
