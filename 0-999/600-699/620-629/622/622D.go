package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var x int64
	// first half
	x = (n - 1) % 2
	if x == 0 {
		x = 2
	}
	for x != n+1 {
		fmt.Fprintf(writer, "%d ", x)
		x += 2
	}
	x -= 2
	for x > 0 {
		fmt.Fprintf(writer, "%d ", x)
		x -= 2
	}
	// middle element
	fmt.Fprintf(writer, "%d ", n)
	// second half
	x = n % 2
	if x == 0 {
		x = 2
	}
	for x != n+2 {
		fmt.Fprintf(writer, "%d ", x)
		x += 2
	}
	x -= 4
	for x > 0 {
		fmt.Fprintf(writer, "%d ", x)
		x -= 2
	}
	fmt.Fprintln(writer)
}
