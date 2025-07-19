package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var a, b int
	if _, err := fmt.Fscan(os.Stdin, &a, &b); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Print increasing sequence
	for i := b + 1; i <= b+a+1; i++ {
		fmt.Fprintf(writer, "%d ", i)
	}
	// Print decreasing sequence
	for i := b; i > 0; i-- {
		fmt.Fprintf(writer, "%d ", i)
	}
}
