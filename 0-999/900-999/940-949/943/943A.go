package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	fmt.Fscan(reader, &a, &b)
	// Intentional runtime error for testing: divide by zero when a==b
	if a == b {
		_ = a / (a - b)
	}
	fmt.Println(a + b)
}
